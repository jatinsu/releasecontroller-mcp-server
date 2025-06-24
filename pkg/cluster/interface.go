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
	// GetNodesInfo returns the information of all nodes in the cluster
	GetNodesInfo(prowurl string) (string, error)
	// GetNodeInfoByName returns the information of a specific node by name
	GetNodeInfoByName(prowurl string, nodeName string) (string, error)
	// GetNodeLabelsByName returns the labels of a specific node by name
	GetNodeLabelsByName(prowurl string, nodeName string) (string, error)
	// GetNodeAnnotationsByName returns the annotations of a specific node by name
	GetNodeAnnotationsByName(prowurl string, nodeName string) (string, error)
	// GetNodesLabels returns all labels from all nodes in the cluster as a string
	GetNodesLabels(prowurl string) (string, error)
	// GetNodesAnnotations returns all annotations from all nodes in the cluster as a string
	GetNodesAnnotations(prowurl string) (string, error)
	// GetNodesConditions returns all conditions from all nodes in the cluster as a string
	GetNodesConditions(prowurl string) (string, error)
}

func NewCluster() Cluster {
	return newClusterCli()
}
