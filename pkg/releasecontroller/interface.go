package releasecontroller

// ReleaseController interface
type ReleaseController interface {
	// ListReleaseControllers lists the available release controllers to use
	ListReleaseControllers() string
	// GetOKDReleaseController returns the OKD release controller URL
	GetOKDReleaseController() string
	// GetOCPReleaseController returns the OpenShift release controller URL
	GetOCPReleaseController() string
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
	// GetJobInfoForRelease gets the build log file for the particular job
	GetJobInfoForRelease(url string) (string, error)
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
