package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	api "github.com/Prashanth684/releasecontroller-mcp-server/pkg/api"
)

// FetchURL fetches data from the given URL and returns it as a string
func FetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// FetchJSONBytes fetches JSON data from the given URL and returns it as a byte slice.
func FetchJSONBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %d %s", resp.StatusCode, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return data, nil
}

// FetchTopLevelKeys fetches JSON from the URL and returns the top-level keys only.
func FetchTopLevelKeys(data []byte) ([]string, error) {
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		return nil, fmt.Errorf("error unmarshaling top-level JSON: %w", err)
	}
	var keys []string
	for key := range top {
		keys = append(keys, key)
	}
	return keys, nil
}

// ParseRelease parses a JSON byte slice into a Release struct
func ParseRelease(data []byte) (*api.Release, error) {
	var r api.Release
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// FilterAcceptedTags filters only tags with Phase == "Accepted"
func FilterAcceptedTags(release *api.Release) []api.Tag {
	var accepted []api.Tag
	for _, tag := range release.Tags {
		if tag.Phase == "Accepted" {
			accepted = append(accepted, tag)
		}
	}
	return accepted
}

// FilterRejectedTags filters only tags with Phase == "Rejected"
func FilterRejectedTags(release *api.Release) []api.Tag {
	var rejected []api.Tag
	for _, tag := range release.Tags {
		if tag.Phase == "Rejected" {
			rejected = append(rejected, tag)
		}
	}
	return rejected
}

// ParseAPIReleaseInfo converts raw JSON bytes into APIReleaseInfo
func ParseAPIReleaseInfo(data []byte) (*api.APIReleaseInfo, error) {
	var info api.APIReleaseInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return &info, nil
}

func ParseVerificationStatusMap(data []byte) (api.VerificationStatusMap, error) {
	var vsm api.VerificationStatusMap
	err := json.Unmarshal(data, &vsm)
	return vsm, err
}

func ParseVerificationJobsSummary(data []byte) (*api.VerificationJobsSummary, error) {
	var summary api.VerificationJobsSummary
	err := json.Unmarshal(data, &summary)
	return &summary, err
}

func ParseUpgradeHistoryList(data []byte) ([]api.UpgradeHistory, error) {
	var upgrades []api.UpgradeHistory
	err := json.Unmarshal(data, &upgrades)
	return upgrades, err
}

func ParseChangeLog(data []byte) (*api.ChangeLog, error) {
	var changelog api.ChangeLog
	err := json.Unmarshal(data, &changelog)
	return &changelog, err
}

func ParseChangeLogReleaseInfo(data []byte) (*api.ChangeLogReleaseInfo, error) {
	var info api.ChangeLogReleaseInfo
	err := json.Unmarshal(data, &info)
	return &info, err
}

func ParseChangeLogComponentInfoList(data []byte) ([]api.ChangeLogComponentInfo, error) {
	var components []api.ChangeLogComponentInfo
	err := json.Unmarshal(data, &components)
	return components, err
}

func ParseChangeLogImageInfoList(data []byte) ([]api.ChangeLogImageInfo, error) {
	var images []api.ChangeLogImageInfo
	err := json.Unmarshal(data, &images)
	return images, err
}

func ParseCommitInfoList(data []byte) ([]api.CommitInfo, error) {
	var commits []api.CommitInfo
	err := json.Unmarshal(data, &commits)
	return commits, err
}

func ExtractProwJobInfo(jobURL string) (string, string, error) {
	u, err := url.Parse(jobURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %w", err)
	}
	parts := strings.Split(u.Path, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("unexpected URL path structure")
	}
	jobID := parts[len(parts)-1]
	jobName := parts[len(parts)-2]
	return jobName, jobID, nil
}

// ExtractTestNameFromURL extracts the first "e2e-*" segment from a prow job URL
func ExtractTestNameFromURL(url string) (string, error) {
	re := regexp.MustCompile(`(?:ocp-)?e2e-[^/]+`)
	match := re.FindString(url)
	if match == "" {
		return "", fmt.Errorf("no e2e test name found in URL: %s", url)
	}
	return strings.TrimPrefix(match, "/"), nil
}

// ExtractStepName parses a log line and extracts the step name
func ExtractStepName(logLine string) (string, error) {
	re := regexp.MustCompile(`Step (.*?) failed after`)
	match := re.FindStringSubmatch(logLine)
	if len(match) < 2 {
		return "", fmt.Errorf("no step name found in line: %s", logLine)
	}
	return strings.TrimSpace(match[1]), nil
}

func CompactTestLogs(input string, threshold float64) string {
	lines := strings.Split(input, "\n")
	var b strings.Builder

	if threshold > 0.6 {
		inBlock := false
		for _, line := range lines {
			if !inBlock && strings.HasPrefix(line, "started:") {
				inBlock = true
			}

			if inBlock {
				if strings.HasPrefix(line, "started:") || strings.HasPrefix(line, "passed: ") || strings.HasPrefix(line, "skipped: ") {
					continue
				}
				b.WriteString(line + "\n")

				//if strings.Contains(line, "Shutting down the monitor") {
				//	break
				//}
			}
		}
	}

	// Couldn't compact any logs, return the original input
	if len(b.String()) <= 0 {
		return DeduplicateLogsWithWindow(input, threshold, 5)
	}
	return DeduplicateLogsWithWindow(b.String(), threshold, 5)
}

func ExtractFailingTestsBlock(input string) string {
	lines := strings.Split(input, "\n")
	var b strings.Builder
	inBlock := false
	for _, line := range lines {
		if strings.Contains(line, "Failing tests:") {
			inBlock = true
		}
		if inBlock {
			// Stop before the end marker
			if strings.Contains(line, "Writing JUnit report to") {
				break
			}
			b.WriteString(line + "\n")
		}
	}
	// If no failing tests block was found, return the original input
	if len(b.String()) <= 0 {
		return DeduplicateLogsWithWindow(input, 1.0, 5)
	}
	return b.String()
}

func ExtractFailedJobsFromAggregate(logData string) map[string]bool {
	jobPIDMap := map[string]string{}
	failedPIDs := map[string]bool{}
	result := map[string]bool{}

	// Match lines like: "********** Starting testcase analysis for: <job>"
	jobRegex := regexp.MustCompile(`\*+ Starting testcase analysis for: (.+)`)
	// Match lines like: "[Tue Jun 10 19:10:22 UTC 2025] <pid> finished with ret=1"
	failRegex := regexp.MustCompile(`\] (\d+) finished with ret=1`)

	lines := strings.Split(logData, "\n")
	var currentJob string
	for _, line := range lines {
		if matches := jobRegex.FindStringSubmatch(line); len(matches) == 2 {
			currentJob = matches[1]
		}
		if strings.HasPrefix(line, "PID is ") && currentJob != "" {
			pid := strings.TrimPrefix(line, "PID is ")
			jobPIDMap[pid] = currentJob
		}
		if matches := failRegex.FindStringSubmatch(line); len(matches) == 2 {
			failedPIDs[matches[1]] = true
		}
	}

	for pid, job := range jobPIDMap {
		result[strings.TrimSpace(job)] = failedPIDs[pid]
	}

	return result
}

func FetchAggregateJobFailures(baseUrl, logData string) (string, error) {
	failJobMap := ExtractFailedJobsFromAggregate(logData)
	if len(failJobMap) == 0 {
		return "", fmt.Errorf("no failed jobs found in the provided log data")
	}
	var builder strings.Builder
	for job, failed := range failJobMap {
		if !failed {
			continue
		}

		data, err := FetchURL(baseUrl + "/" + job + "/" + job + ".log")
		if err != nil {
			return "", fmt.Errorf("error fetching URL %q: %w", baseUrl, err)
		}
		// Find the index of the "summary:" line
		idx := strings.Index(data, "summary:")
		if idx < 0 {
			return "", fmt.Errorf("no summary: found in log for job %q", job)
		}
		builder.WriteString(data[idx:])
		builder.WriteString("\n---\n\n")
	}
	return builder.String(), nil
}

func GetGatherExtraFolderPath(prowurl string) (string, error) {
	name, id, err := ExtractProwJobInfo(prowurl)
	if err != nil {
		return "", fmt.Errorf("error extracting job info: %w", err)
	}
	testName, err := ExtractTestNameFromURL(prowurl)
	if err != nil {
		return "", fmt.Errorf("error fetching test name: %w", err)
	}
	gatherExtraURL := fmt.Sprintf("https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/logs/%s/%s/artifacts/%s/gather-extra/artifacts/", name, id, testName)
	return gatherExtraURL, nil
}

func IndentMultiline(s, indent string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}
