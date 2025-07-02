package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Types originally from origin monitorapi package
type Locator struct {
	Type string            `json:"type"`
	Keys map[string]string `json:"keys"`
}

type Message struct {
	Reason       string            `json:"reason"`
	Cause        string            `json:"cause"`
	HumanMessage string            `json:"humanMessage"`
	Annotations  map[string]string `json:"annotations"`
}
type EventInterval struct {
	Level             string  `json:"level"`
	Display           bool    `json:"display"`
	Source            string  `json:"source,omitempty"`
	StructuredLocator Locator `json:"locator"`
	StructuredMessage Message `json:"message"`

	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
	// Filename is the base filename we read the intervals from in gcs. If multiple,
	// that usually means one for upgrade and one for conformance portions of the job run.
	// TODO: this may need to be revisited once we're further along with the UI/new schema.
	Filename string `json:"filename"`
}

type Report struct {
	Items []EventInterval `json:"items"`
}

func GetSpyglassFileNames(logsPath, testName, stepFolder string) ([]string, error) {
	// Compile the regex pattern
	pattern := `e2e-timelines_spyglass_.*\.json$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Make the HTTP GET request
	artifactURL := fmt.Sprintf("https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/%s/artifacts/%s/%s/artifacts/junit/", logsPath, testName, stepFolder)
	resp, err := http.Get(artifactURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Collect matching file names
	var matches []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if re.MatchString(text) {
			matches = append(matches, text)
		}
	})

	return matches, nil
}

func GetErrorAndWarningFromSpyglassFile(spyglassFilePath string) (string, error) {
	resp, err := http.Get(spyglassFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to fetch spyglass file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var events Report
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return "", fmt.Errorf("failed to decode spyglass file: %w %s", err, spyglassFilePath)
	}

	var result []EventInterval
	for _, event := range events.Items {
		if event.Level == "Error" || event.Level == "Warning" {
			result = append(result, event)
		}
	}
	if len(result) == 0 {
		return "no errors and warnings!", nil
	}
	//convert event to JSON string
	resultStr, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal events: %w", err)
	}
	return string(resultStr), nil
}

// Given a test name check for the Locator.Keys objects in the spyglass data to see if there is an entry with key "e2e-test" that matches the test name. return result as string
func GetSpyglassDataRelevantToTestFailure(spyglassFilePath, testName string) (string, error) {
	errorEvents, err := GetErrorAndWarningFromSpyglassFile(spyglassFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get error and warning events: %w", err)
	}
	if errorEvents == "" {
		return "no spyglass data", nil // No errors or warnings found
	}
	var errorEventObj []EventInterval
	err = json.Unmarshal([]byte(errorEvents), &errorEventObj)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal error events: %w %s", err, spyglassFilePath)
	}
	testFound := false
	var relevantEvents string
	for _, event := range errorEventObj {
		//Extract only the "Source" field and the message.HumanMessage field and construct it as a string and append it to relevantEvents
		if event.Source == "" || event.StructuredMessage.HumanMessage == "" {
			continue // Skip events without Source or HumanMessage
		}
		if event.StructuredLocator.Keys == nil {
			continue // Skip events without keys
		}
		if event.From == nil || event.To == nil {
			continue // Skip events without a time range
		}
		testLocator := event.StructuredLocator.Keys["e2e-test"]
		if testLocator == "" {
			testLocator = "Not a test"
		}
		// Construct the relevant event string
		eventString := fmt.Sprintf("Source: %s Type: %s Locator: 'test: %s' Reason: %s HumanMessage: %s From: %s To: %s\n",
			event.Source,
			event.StructuredLocator.Type,
			testLocator,
			event.StructuredMessage.Reason,
			event.StructuredMessage.HumanMessage,
			event.From.Format(time.RFC3339),
			event.To.Format(time.RFC3339),
		)

		relevantEvents += eventString
		for _, key := range event.StructuredLocator.Keys {
			if key == "e2e-test" && event.StructuredLocator.Keys["e2e-test"] == testName {
				testFound = true
				break
			}
		}
		if testFound {
			break // Stop searching once we find the test
		}
	}
	if relevantEvents == "" {
		return "no relevant data", nil // No relevant events found
	}
	return relevantEvents, nil
}
