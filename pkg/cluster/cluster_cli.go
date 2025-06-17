package cluster

import (
	"fmt"
	"strings"

	"github.com/Prashanth684/releasecontroller-mcp-server/pkg/utils"
	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterCli struct {
}

func (c *clusterCli) GetPodsInState(prowurl string, state string) (string, error) {
	//Fetch the url of the extra folder
	artifactURL, err := utils.GetGatherExtraFolderPath(prowurl)
	if err != nil {
		return "", fmt.Errorf("error getting gather extra folder path: %w", err)
	}
	// Download the pods.json file from the artifact URL
	pods, err := utils.LoadPodsFromFile(artifactURL + "pods.json")
	if err != nil {
		return "", fmt.Errorf("error loading pods: %w", err)
	}
	var result string
	switch state {
	case "CrashLoopBackOff":
		result = utils.CrashLoopBackOffSummary(pods)
	case "Pending":
		result = utils.PendingPodsSummary(pods)
	case "Init":
		result = utils.InitStateSummary(pods)
	case "Error":
		result = utils.ErrorStateSummary(pods)
	case "Running":
		result = utils.RunningPodsSummary(pods)
	default:
		result = utils.AllPodsSummary(pods)
	}
	return result, nil
}

func (c *clusterCli) GetPodsInNamespace(prowurl string, namespace string) (string, error) {
	//Fetch the url of the extra folder
	artifactURL, err := utils.GetGatherExtraFolderPath(prowurl)
	if err != nil {
		return "", fmt.Errorf("error getting gather extra folder path: %w", err)
	}
	// Download the pods.json file from the artifact URL
	pods, err := utils.LoadPodsFromFile(artifactURL + "pods.json")
	if err != nil {
		return "", fmt.Errorf("error loading pods: %w", err)
	}
	return utils.FilterPodsByNamespaceAsString(pods, namespace), nil
}

func (c *clusterCli) GetPodsInNode(prowurl string, nodeName string) (string, error) {
	//Fetch the url of the extra folder
	artifactURL, err := utils.GetGatherExtraFolderPath(prowurl)
	if err != nil {
		return "", fmt.Errorf("error getting gather extra folder path: %w", err)
	}
	// Download the pods.json file from the artifact URL
	pods, err := utils.LoadPodsFromFile(artifactURL + "pods.json")
	if err != nil {
		return "", fmt.Errorf("error loading pods: %w", err)
	}
	return utils.FilterPodsByNodeAsString(pods, nodeName), nil
}

// Extract status summary (Available, Progressing, Degraded) for each operator
func (c *clusterCli) GetClusterOperatorStatusSummary(prowurl string) (string, error) {
	//Fetch the url of the extra folder
	artifactURL, err := utils.GetGatherExtraFolderPath(prowurl)
	if err != nil {
		return "", fmt.Errorf("error getting gather extra folder path: %w", err)
	}
	// Download the clusteroperators.json file from the artifact URL
	operators, err := utils.LoadClusterOperatorsFromFile(artifactURL + "clusteroperators.json")
	if err != nil {
		return "", fmt.Errorf("error loading cluster operators: %w", err)
	}
	var b strings.Builder
	for _, op := range operators {
		name := op.Name
		var available, progressing, degraded string

		for _, cond := range op.Status.Conditions {
			switch cond.Type {
			case configv1.OperatorAvailable:
				available = fmt.Sprintf("%s (Reason: %s)", cond.Status, cond.Reason)
			case configv1.OperatorProgressing:
				progressing = fmt.Sprintf("%s (Reason: %s)", cond.Status, cond.Reason)
			case configv1.OperatorDegraded:
				degraded = fmt.Sprintf("%s (Reason: %s)", cond.Status, cond.Reason)
			}
		}

		fmt.Fprintf(&b, "Operator: %s\n  Available: %s\n  Progressing: %s\n  Degraded: %s\n\n",
			name, available, progressing, degraded)
	}
	return b.String(), nil
}

func (c *clusterCli) GetClusterVersionSummary(prowurl string) (string, error) {
	// Fetch the url of the extra folder
	artifactURL, err := utils.GetGatherExtraFolderPath(prowurl)
	if err != nil {
		return "", fmt.Errorf("error getting gather extra folder path: %w", err)
	}
	// Download the clusterversion.json file from the artifact URL
	clusterVersion, err := utils.LoadClusterVersionFromFile(artifactURL + "clusterversion.json")
	if err != nil {
		return "", fmt.Errorf("error loading cluster version: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Cluster Version: %s\n", clusterVersion.Status.Desired.Version)
	fmt.Fprintf(&b, "Desired Image: %s\n", clusterVersion.Status.Desired.Image)
	fmt.Fprintf(&b, "Desired URL: %s\n", clusterVersion.Status.Desired.URL)
	fmt.Fprintf(&b, "Available Updates: %d\n", len(clusterVersion.Status.AvailableUpdates))
	for _, update := range clusterVersion.Status.AvailableUpdates {
		fmt.Fprintf(&b, "  - %s\n", update.Version)
	}

	fmt.Fprintf(&b, "\nUpdate History:\n")
	for _, hist := range clusterVersion.Status.History {
		fmt.Fprintf(&b, "  - Version: %s | State: %s | Verified: %t\n    Image: %s\n",
			hist.Version, hist.State, hist.Verified, hist.Image)
		fmt.Fprintf(&b, "    Started: %s | Completed: %s\n",
			hist.StartedTime.Time.Format("2006-01-02 15:04:05"),
			func(t *metav1.Time) string {
				if t == nil {
					return "N/A"
				}
				return t.Time.Format("2006-01-02 15:04:05")
			}(hist.CompletionTime))
		if hist.AcceptedRisks != "" {
			fmt.Fprintf(&b, "    Accepted Risks:\n    %s\n", utils.IndentMultiline(hist.AcceptedRisks, "    "))
		}
	}
	return b.String(), nil
}

func newClusterCli() *clusterCli {
	return &clusterCli{}
}
