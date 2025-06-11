package releasecontroller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	OKDReleaseController = "amd64.origin.releases.ci.openshift.org"
	OCPReleaseController = "amd64.ocp.releases.ci.openshift.org"
)

type releaseControllerCli struct {
	releaseControllers []string
}

// Tag represents a single entry in the "tags" array
type Tag struct {
	Name        string `json:"name"`
	Phase       string `json:"phase"`
	PullSpec    string `json:"pullSpec"`
	DownloadURL string `json:"downloadURL"`
}

// Release represents the full JSON structure
type Release struct {
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

// Copied from release controller API: https://github.com/openshift/release-controller/blob/main/pkg/release-controller/types.go
// APIReleaseInfo encapsulates the release verification results and upgrade history for a release tag.
type APIReleaseInfo struct {
	// Name is the name of the release tag.
	Name string `json:"name"`
	// Phase is the phase of the release tag.
	Phase string `json:"phase"`
	// Results is the status of the release verification jobs for this release tag
	Results *VerificationJobsSummary `json:"results,omitempty"`
	// UpgradesTo is the list of UpgradeHistory "to" this release tag
	UpgradesTo []UpgradeHistory `json:"upgradesTo,omitempty"`
	//UpgradesFrom is the list of UpgradeHistory "from" this release tag
	UpgradesFrom []UpgradeHistory `json:"upgradesFrom,omitempty"`
	//ChangeLog is the html representation of the changes included in this release tag
	ChangeLog []byte `json:"changeLog,omitempty"`
	//ChangeLogJson is the json representation of the changes included in this release tag
	ChangeLogJson ChangeLog `json:"changeLogJson,omitempty"`
}

type VerificationStatus struct {
	State   string `json:"state"`
	URL     string `json:"url"`
	Retries int    `json:"retries,omitempty"`
	// TransitionTime *metav1.Time `json:"transitionTime,omitempty"`
}

type VerificationStatusMap map[string]*VerificationStatus

// VerificationJobsSummary an organized, by job type, collection of VerificationStatusMap objects
type VerificationJobsSummary struct {
	BlockingJobs  VerificationStatusMap `json:"blockingJobs,omitempty"`
	InformingJobs VerificationStatusMap `json:"informingJobs,omitempty"`
	PendingJobs   VerificationStatusMap `json:"pendingJobs,omitempty"`
}

type UpgradeResult struct {
	State string `json:"state"`
	URL   string `json:"url"`
}

type UpgradeHistory struct {
	From string
	To   string

	Success int
	Failure int
	Total   int

	History map[string]UpgradeResult
}

// ChangeLog represents the data structure that oc returns when providing a changelog in JSON format
// TODO: This is being carried from changes in openshift/oc.  These changes should be removed if/when we bump up our k8s dependencies up to the latest/greatest version.  We're currently pinned at: v0.24.2
type ChangeLog struct {
	From ChangeLogReleaseInfo `json:"from"`
	To   ChangeLogReleaseInfo `json:"to"`

	Components    []ChangeLogComponentInfo `json:"components,omitempty"`
	NewImages     []ChangeLogImageInfo     `json:"newImages,omitempty"`
	RemovedImages []ChangeLogImageInfo     `json:"removedImages,omitempty"`
	RebuiltImages []ChangeLogImageInfo     `json:"rebuiltImages,omitempty"`
	UpdatedImages []ChangeLogImageInfo     `json:"updatedImages,omitempty"`
}

type ChangeLogReleaseInfo struct {
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	//	Digest       digest.Digest `json:"digest"`
	PromotedFrom string `json:"promotedFrom,omitempty"`
}

type ChangeLogComponentInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	VersionUrl string `json:"versionUrl,omitempty"`
	From       string `json:"from,omitempty"`
	FromUrl    string `json:"fromUrl,omitempty"`
	DiffUrl    string `json:"diffUrl,omitempty"`
}

type ChangeLogImageInfo struct {
	Name          string       `json:"name"`
	Path          string       `json:"path"`
	ShortCommit   string       `json:"shortCommit,omitempty"`
	Commit        string       `json:"commit,omitempty"`
	ImageRef      string       `json:"imageRef,omitempty"`
	Commits       []CommitInfo `json:"commits,omitempty"`
	FullChangeLog string       `json:"fullChangeLog,omitempty"`
}

type CommitInfo struct {
	Bugs      map[string]string `json:"bugs,omitempty"`
	Issues    map[string]string `json:"issues,omitempty"`
	Subject   string            `json:"subject,omitempty"`
	PullID    int               `json:"pullID,omitempty"`
	PullURL   string            `json:"pullURL,omitempty"`
	CommitID  string            `json:"commitID,omitempty"`
	CommitURL string            `json:"commitURL,omitempty"`
}

// End of releaseController API

// FetchURL fetches data from the given URL and returns it as a string
func fetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// fetchJSONBytes fetches JSON data from the given URL and returns it as a byte slice.
func fetchJSONBytes(url string) ([]byte, error) {
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

// fetchTopLevelKeys fetches JSON from the URL and returns the top-level keys only.
func fetchTopLevelKeys(data []byte) ([]string, error) {
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

// parseRelease parses a JSON byte slice into a Release struct
func parseRelease(data []byte) (*Release, error) {
	var r Release
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// filterAcceptedTags filters only tags with Phase == "Accepted"
func filterAcceptedTags(release *Release) []Tag {
	var accepted []Tag
	for _, tag := range release.Tags {
		if tag.Phase == "Accepted" {
			accepted = append(accepted, tag)
		}
	}
	return accepted
}

// filterRejectedTags filters only tags with Phase == "Accepted"
func filterRejectedTags(release *Release) []Tag {
	var accepted []Tag
	for _, tag := range release.Tags {
		if tag.Phase == "Rejected" {
			accepted = append(accepted, tag)
		}
	}
	return accepted
}

// ParseAPIReleaseInfo converts raw JSON bytes into APIReleaseInfo
func parseAPIReleaseInfo(data []byte) (*APIReleaseInfo, error) {
	var info APIReleaseInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	return &info, nil
}

func parseVerificationStatusMap(data []byte) (VerificationStatusMap, error) {
	var vsm VerificationStatusMap
	err := json.Unmarshal(data, &vsm)
	return vsm, err
}

func parseVerificationJobsSummary(data []byte) (*VerificationJobsSummary, error) {
	var summary VerificationJobsSummary
	err := json.Unmarshal(data, &summary)
	return &summary, err
}

func parseUpgradeHistoryList(data []byte) ([]UpgradeHistory, error) {
	var upgrades []UpgradeHistory
	err := json.Unmarshal(data, &upgrades)
	return upgrades, err
}

func parseChangeLog(data []byte) (*ChangeLog, error) {
	var changelog ChangeLog
	err := json.Unmarshal(data, &changelog)
	return &changelog, err
}

func parseChangeLogReleaseInfo(data []byte) (*ChangeLogReleaseInfo, error) {
	var info ChangeLogReleaseInfo
	err := json.Unmarshal(data, &info)
	return &info, err
}

func parseChangeLogComponentInfoList(data []byte) ([]ChangeLogComponentInfo, error) {
	var components []ChangeLogComponentInfo
	err := json.Unmarshal(data, &components)
	return components, err
}

func parseChangeLogImageInfoList(data []byte) ([]ChangeLogImageInfo, error) {
	var images []ChangeLogImageInfo
	err := json.Unmarshal(data, &images)
	return images, err
}

func parseCommitInfoList(data []byte) ([]CommitInfo, error) {
	var commits []CommitInfo
	err := json.Unmarshal(data, &commits)
	return commits, err
}

func extractProwJobInfo(jobURL string) (string, string, error) {
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

// extractTestNameFromURL extracts the first "e2e-*" segment from a prow job URL
func extractTestNameFromURL(url string) (string, error) {
	re := regexp.MustCompile(`e2e-[^/]+`)
	match := re.FindString(url)

	if match == "" {
		return "", fmt.Errorf("no e2e test name found in URL: %s", url)
	}
	return strings.TrimPrefix(match, "/"), nil
}

// extractStepName parses a log line and extracts the step name
func extractStepName(logLine string) (string, error) {
	re := regexp.MustCompile(`Step (.*?) failed after`)
	match := re.FindStringSubmatch(logLine)

	if len(match) < 2 {
		return "", fmt.Errorf("no step name found in line: %s", logLine)
	}
	return strings.TrimSpace(match[1]), nil
}

// ListReleaseControllers lists the available release controllers to use
func (r *releaseControllerCli) ListReleaseControllers() string {
	return strings.Join(r.releaseControllers, ",")
}

// GetOKDReleaseController returns the OKD release controller host
func (r *releaseControllerCli) GetOKDReleaseController() string {
	return OKDReleaseController
}

// GetOCPReleaseController returns the OCP release controller host
func (r *releaseControllerCli) GetOCPReleaseController() string {
	return OCPReleaseController
}

// ListReleaseStreams lists all the releases from all the streams in the release controller
func (r *releaseControllerCli) ListReleaseStreams(releasecontroller string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestreams/all", releasecontroller))
	if err != nil {
		return "", fmt.Errorf("error fetching release streams: %w", err)
	}
	topKeys, err := fetchTopLevelKeys(data)
	if err != nil {
		return "", fmt.Errorf("error fetching top-level keys: %w", err)
	}
	return strings.Join(topKeys, ", "), nil
}

// LatestAcceptedRelease gets the latest accepted release for a given stream
func (r *releaseControllerCli) LatestAcceptedRelease(releasecontroller, stream string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/tags", releasecontroller, stream))
	if err != nil {
		return "", fmt.Errorf("error fetching release tags: %w", err)
	}
	release, err := parseRelease(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release data: %w", err)
	}
	acceptedTags := filterAcceptedTags(release)
	if len(acceptedTags) == 0 {
		return "", errors.New("no accepted tags found")
	}
	var latestAcceptedTag string
	for _, tag := range acceptedTags {
		if latestAcceptedTag == "" || tag.Name > latestAcceptedTag {
			latestAcceptedTag = tag.Name
		}
	}
	return latestAcceptedTag, nil
}

// LatestRejectedRelease gets the latest rejected release for a given stream
func (r *releaseControllerCli) LatestRejectedRelease(releasecontroller, stream string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/tags", releasecontroller, stream))
	if err != nil {
		return "", fmt.Errorf("error fetching release tags: %w", err)
	}
	release, err := parseRelease(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release data: %w", err)
	}
	rejectedTags := filterRejectedTags(release)
	if len(rejectedTags) == 0 {
		return "", errors.New("no rejected tags found")
	}
	var latestRejectedTag string
	for _, tag := range rejectedTags {
		if latestRejectedTag == "" || tag.Name > latestRejectedTag {
			latestRejectedTag = tag.Name
		}
	}
	return latestRejectedTag, nil
}

// ListFailedJobsInRelease lists all the failed jobs in a given release
func (r *releaseControllerCli) ListFailedJobsInRelease(releasecontroller, stream, tag string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := parseAPIReleaseInfo(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release info: %w", err)
	}
	var failedJobs []string
	for jobName, status := range info.Results.BlockingJobs {
		if status.State == "Failed" {
			failedJobs = append(failedJobs, fmt.Sprintf("%s: %s", jobName, status.URL))
		}
	}
	for jobName, status := range info.Results.InformingJobs {
		if status.State == "Failed" {
			failedJobs = append(failedJobs, fmt.Sprintf("%s: %s", jobName, status.URL))
		}
	}
	if len(failedJobs) == 0 {
		return "No failed jobs found", nil
	}
	return strings.Join(failedJobs, "\n"), nil
}

// ListComponentsInRelease lists the kubectl, kubernetes, coreos and tests versions in the release
func (r *releaseControllerCli) ListComponentsInRelease(releasecontroller, stream, tag string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := parseAPIReleaseInfo(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release info: %w", err)
	}
	var components []string
	for _, component := range info.ChangeLogJson.Components {
		components = append(components, fmt.Sprintf("%s: %s", component.Name, component.Version))
	}
	if len(components) == 0 {
		return "No components found", nil
	}
	return strings.Join(components, "\n"), nil
}

// GetJobInfoForRelease gets the build log file for the particular job
func (r *releaseControllerCli) GetJobInfoForRelease(prowurl string) (string, error) {
	name, id, err := extractProwJobInfo(prowurl)
	if err != nil {
		return "", fmt.Errorf("error extracting job info: %w", err)
	}
	joburl := fmt.Sprintf("https://storage.googleapis.com/test-platform-results/logs/%s/%s/build-log.txt", name, id)
	data, err := fetchURL(joburl)
	if err != nil {
		return "", fmt.Errorf("error fetching job log: %w", err)
	}
	stepName, err := extractStepName(data)
	if err != nil {
		return data, nil
	}
	testName, err := extractTestNameFromURL(prowurl)
	if err != nil {
		return "", fmt.Errorf("error fetching test name: %w", err)
	}
	if !strings.HasPrefix(stepName, testName+"-") {
		return "", fmt.Errorf("stepName does not start with testName prefix")
	}
	stepFolder := strings.TrimPrefix(stepName, testName+"-")
	artifactURL := fmt.Sprintf("https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/logs/%s/%s/artifacts/%s/%s/build-log.txt", name, id, testName, stepFolder)
	testLogs, err := fetchURL(artifactURL)
	if err != nil {
		return "", fmt.Errorf("error fetching test logs: %w", err)
	}
	const marker = "Failing tests:"
	idx := strings.Index(testLogs, marker)
	if idx == -1 {
		return testLogs, nil // No marker found, return the full log
	}
	// Slice after the marker
	return data[idx+len(marker):], nil
}

// List issues which are features from updated images commits - excludes OCPBUGS/CVEs
func (r *releaseControllerCli) ListFeaturesFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := parseAPIReleaseInfo(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release info: %w", err)
	}
	var components []string
	for _, component := range info.ChangeLogJson.UpdatedImages {
		for _, commit := range component.Commits {
			if len(commit.Issues) > 0 {
				for issue, url := range commit.Issues {
					if strings.HasPrefix(issue, "OCPBUGS-") || strings.HasPrefix(issue, "CVE-") {
						continue // Skip OCPBUGS and CVEs
					}
					// Add the issue to the components list
					// Format: "issue: url (component name)"
					components = append(components, fmt.Sprintf("%s: %s", issue, url))
				}
			}
		}
	}
	if len(components) == 0 {
		return "No issues found in updated images commits", nil
	}
	return strings.Join(components, "\n"), nil
}

// List issues which are bugs from updated images commits
func (r *releaseControllerCli) ListBugsFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := parseAPIReleaseInfo(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release info: %w", err)
	}
	var components []string
	for _, component := range info.ChangeLogJson.UpdatedImages {
		for _, commit := range component.Commits {
			if len(commit.Issues) > 0 {
				for issue, url := range commit.Issues {
					if strings.HasPrefix(issue, "OCPBUGS-") {
						// Add the issue to the components list
						// Format: "issue: url (component name)"
						components = append(components, fmt.Sprintf("%s: %s", issue, url))
					}
				}
			}
		}
	}
	if len(components) == 0 {
		return "No issues found in updated images commits", nil
	}
	return strings.Join(components, "\n"), nil
}

// List issues which are CVEs from updated images commits
func (r *releaseControllerCli) ListCVEsFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error) {
	data, err := fetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := parseAPIReleaseInfo(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release info: %w", err)
	}
	var components []string
	for _, component := range info.ChangeLogJson.UpdatedImages {
		for _, commit := range component.Commits {
			if len(commit.Issues) > 0 && strings.Contains(commit.Subject, "CVE") {
				// Only consider commits with CVE issues
				for issue, url := range commit.Issues {
					// Add the issue to the components list
					// Format: "issue: url (component name)"
					components = append(components, fmt.Sprintf("%s: %s", issue, url))
				}
			}
		}
	}
	if len(components) == 0 {
		return "No issues found in updated images commits", nil
	}
	return strings.Join(components, "\n"), nil
}

func newReleaseControllerCli() *releaseControllerCli {
	return &releaseControllerCli{
		releaseControllers: []string{OKDReleaseController, OCPReleaseController},
	}
}
