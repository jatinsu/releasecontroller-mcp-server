package releasecontroller

// ReleaseController interface
type ReleaseController interface {
	// ListReleaseControllers lists the available release controllers to use
	ListReleaseControllers() string
	// GetOKDReleaseController returns the OKD release controller URL
	GetOKDReleaseController() string
	// GetOCPReleaseController returns the OpenShift release controller URL
	GetOCPReleaseController() string
	//GetMultiReleaseController returns the multi-arch release controller URL
	GetMultiReleaseController() string
	//GetARM64ReleaseController returns the ARM64 release controller URL
	GetARM64ReleaseController() string
	// GetPPC64LEReleaseController returns the PPC64LE release controller URL
	GetPPC64LEReleaseController() string
	// GetS390XReleaseController returns the S390X release controller URL
	GetS390XReleaseController() string
	// ListReleaseStreams lists all the release streams in the release controller
	ListReleaseStreams(releasecontroller string) (string, error)
	// LatestAcceptedRelease gets the latest accepted release for a given stream
	LatestAcceptedRelease(releasecontroller, stream string) (string, error)
	// LatestRejectedRelease gets the latest rejected release for a given stream
	LatestRejectedRelease(releasecontroller, stream string) (string, error)
	// ListFailedJobsInRelease lists all the failed jobs in a given release
	ListFailedJobsInRelease(releasecontroller, stream, tag string) (string, error)
	// ListComponentsInRelease lists the kubectl, kubernetes, coreos and tests versions in the release
	ListComponentsInRelease(releasecontroller, stream, tag string) (string, error)
	// ListTestFailuresForRelease gets the failing tests for the particular job
	ListTestFailuresForRelease(prowurl string) (string, error)
	//GetFlakyTestsForRelease gets the flaky tests for the particular job
	GetFlakyTestsForRelease(prowurl string) (string, error)
	// GetRiskAnalysisData gets the risk analysis data for the particular job
	GetRiskAnalysisData(prowurl string) (string, error)
	// GetSpyglassDataRelevantToTestFailure gets the spyglass data relevant to a test failure
	GetSpyglassDataRelevantToTestFailure(prowurl string, testName string) (string, error)
	//GetTopLevelBuildLog gets the top-level build log for a given Prow job URL
	GetTopLevelBuildLog(prowurl string, LogCompactionThreshold string) (string, error)
	// AnalyzeJobFailuresForRelease gets the build log file for the particular job
	AnalyzeJobFailuresForRelease(url string, LogCompactionThreshold string) (string, error)
	// List issues which are features from updated images commits - excludes OCPBUGS/CVEs
	ListFeaturesFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error)
	// List issues which are bugs from updated images commits
	ListBugsFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error)
	// List issues which are CVEs from updated images commits
	ListCVEsFromUpdatedImagesCommits(releasecontroller, stream, tag string) (string, error)
}

func NewReleaseController() ReleaseController {
	return newReleaseControllerCli()
}
