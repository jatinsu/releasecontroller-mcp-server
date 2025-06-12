package releasecontroller

import (
	"errors"
	"fmt"
	"strings"

	utils "github.com/Prashanth684/releasecontroller-mcp-server/pkg/utils"
)

const (
	OKDReleaseController = "amd64.origin.releases.ci.openshift.org"
	OCPReleaseController = "amd64.ocp.releases.ci.openshift.org"
)

type releaseControllerCli struct {
	releaseControllers []string
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
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestreams/all", releasecontroller))
	if err != nil {
		return "", fmt.Errorf("error fetching release streams: %w", err)
	}
	topKeys, err := utils.FetchTopLevelKeys(data)
	if err != nil {
		return "", fmt.Errorf("error fetching top-level keys: %w", err)
	}
	return strings.Join(topKeys, ", "), nil
}

// LatestAcceptedRelease gets the latest accepted release for a given stream
func (r *releaseControllerCli) LatestAcceptedRelease(releasecontroller, stream string) (string, error) {
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/tags", releasecontroller, stream))
	if err != nil {
		return "", fmt.Errorf("error fetching release tags: %w", err)
	}
	release, err := utils.ParseRelease(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release data: %w", err)
	}
	acceptedTags := utils.FilterAcceptedTags(release)
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
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/tags", releasecontroller, stream))
	if err != nil {
		return "", fmt.Errorf("error fetching release tags: %w", err)
	}
	release, err := utils.ParseRelease(data)
	if err != nil {
		return "", fmt.Errorf("error parsing release data: %w", err)
	}
	rejectedTags := utils.FilterRejectedTags(release)
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
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := utils.ParseAPIReleaseInfo(data)
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
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := utils.ParseAPIReleaseInfo(data)
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

// ListTestFailuresForRelease gets the failing tests for the particular job
func (r *releaseControllerCli) ListTestFailuresForRelease(prowurl string) (string, error) {
	name, id, err := utils.ExtractProwJobInfo(prowurl)
	if err != nil {
		return "", fmt.Errorf("error extracting job info: %w", err)
	}
	joburl := fmt.Sprintf("https://storage.googleapis.com/test-platform-results/logs/%s/%s/build-log.txt", name, id)
	data, err := utils.FetchURL(joburl)
	if err != nil {
		return "", fmt.Errorf("error fetching job log: %w", err)
	}
	stepName, err := utils.ExtractStepName(data)
	if err != nil {
		return "", fmt.Errorf("Could not find failure step - not a test run", err)
	}
	testName, err := utils.ExtractTestNameFromURL(prowurl)
	if err != nil {
		return "", fmt.Errorf("error fetching test name: %w", err)
	}
	if !strings.HasPrefix(stepName, testName+"-") {
		return "", fmt.Errorf("stepName does not start with testName prefix")
	}
	stepFolder := strings.TrimPrefix(stepName, testName+"-")
	artifactURL := fmt.Sprintf("https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/logs/%s/%s/artifacts/%s/%s/build-log.txt", name, id, testName, stepFolder)
	testLogs, err := utils.FetchURL(artifactURL)
	if err != nil {
		return "", fmt.Errorf("error fetching test logs: %w", err)
	}
	return utils.ExtractFailingTestsBlock(testLogs), nil
}

// AnalyzeJobFailuresForRelease gets the build log file for the particular job
func (r *releaseControllerCli) AnalyzeJobFailuresForRelease(prowurl string) (string, error) {
	name, id, err := utils.ExtractProwJobInfo(prowurl)
	if err != nil {
		return "", fmt.Errorf("error extracting job info: %w", err)
	}
	joburl := fmt.Sprintf("https://storage.googleapis.com/test-platform-results/logs/%s/%s/build-log.txt", name, id)
	data, err := utils.FetchURL(joburl)
	if err != nil {
		return "", fmt.Errorf("error fetching job log: %w", err)
	}
	stepName, err := utils.ExtractStepName(data)
	if err != nil {
		return data, nil
	}
	testName, err := utils.ExtractTestNameFromURL(prowurl)
	if err != nil {
		return "", fmt.Errorf("error fetching test name: %w", err)
	}
	if !strings.HasPrefix(stepName, testName+"-") {
		return "", fmt.Errorf("stepName does not start with testName prefix")
	}
	stepFolder := strings.TrimPrefix(stepName, testName+"-")
	artifactURL := fmt.Sprintf("https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/logs/%s/%s/artifacts/%s/%s/build-log.txt", name, id, testName, stepFolder)
	testLogs, err := utils.FetchURL(artifactURL)
	if err != nil {
		return "", fmt.Errorf("error fetching test logs: %w", err)
	}
	testLogs = utils.CompactTestLogs(testLogs)
	return testLogs, nil
}

// List issues which are features from updated images commits - excludes OCPBUGS/CVEs
func (r *releaseControllerCli) ListFeaturesFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error) {
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := utils.ParseAPIReleaseInfo(data)
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
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := utils.ParseAPIReleaseInfo(data)
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
	data, err := utils.FetchJSONBytes(fmt.Sprintf("https://%s/api/v1/releasestream/%s/release/%s", releasecontroller, stream, tag))
	if err != nil {
		return "", fmt.Errorf("error fetching release info: %w", err)
	}
	info, err := utils.ParseAPIReleaseInfo(data)
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
