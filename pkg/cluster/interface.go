package cluster

// ReleaseController interface
type Cluster interface {
	// GetPodsInState returns the pods in a specific state
	GetPodsInState(prowurl string, state string) (string, error)
	// GetPodsInNamespace returns the pods in a specific namespace
	GetPodsInNamespace(prowurl string, namespace string) (string, error)
	// GetPodsInNode returns the pods in a specific node
	GetPodsInNode(prowurl string, nodeName string) (string, error)
	// GetContainersInPod returns the containers in a specific pod
	GetContainersInPod(prowurl string, podName string, namespace string) (string, error)
	// GetContainerLogs returns the logs of a specific container in a pod
	GetContainerLogs(prowurl string, podName string, namespace string, containerName string) (string, error)
	// GetClusterOperatorStatusSummary returns the status summary of cluster operators
	GetClusterOperatorStatusSummary(prowurl string) (string, error)
	// GetClusterVersionSummary returns the cluster version summary
	GetClusterVersionSummary(prowurl string) (string, error)
	// GetClusterNodes returns the list of cluster nodes
	//GetClusterNodes(prowurl string) ([]corev1.Node, error)
}

func NewCluster() Cluster {
	return newClusterCli()
}
