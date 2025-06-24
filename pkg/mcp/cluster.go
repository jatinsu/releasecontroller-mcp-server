package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register the CLI tools for the release controller.
func (s *Server) initCluster() []server.ServerTool {
	return []server.ServerTool{
		{mcp.NewTool("get_pods_in_state",
			mcp.WithDescription("Get pods in a specific state mentioned by the user. The state can be one of: CrashLoopBackOff, Pending, Init, Error, Running, or All."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("state", mcp.Description("State of the pods to filter"), mcp.Required()),
		), func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			state := ctr.Params.Arguments["state"].(string)
			result, err := s.cluster.GetPodsInState(prowurl, state)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_cluster_operator_status_summary",
			mcp.WithDescription("Get status summary of cluster operators. Clearly list the available, progressing, and degraded states of each operator. Format the output neatly with operator name, available, progressing, and degraded states."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetClusterOperatorStatusSummary(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_cluster_version_summary",
			mcp.WithDescription("Get the cluster version summary including the current version, desired version, and available updates. Format the output neatly."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetClusterVersionSummary(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_pods_in_namespace",
			mcp.WithDescription("Get pods in a specific namespace. Format the output neatly with the pod name and namespace."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace to filter pods"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			namespace := ctr.Params.Arguments["namespace"].(string)
			result, err := s.cluster.GetPodsInNamespace(prowurl, namespace)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_pods_in_node",
			mcp.WithDescription("Get pods in a specific node. Format the output neatly with the pod name, namespace, and node name."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("nodeName", mcp.Description("Node name to filter pods"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			nodeName := ctr.Params.Arguments["nodeName"].(string)
			result, err := s.cluster.GetPodsInNode(prowurl, nodeName)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_containers_in_pod",
			mcp.WithDescription("Get containers in a specific pod. Format the ouput neatly with the pod name, namespace, and node name."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("podName", mcp.Description("Pod name to filter containers"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace of the pod"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			podName := ctr.Params.Arguments["podName"].(string)
			namespace := ctr.Params.Arguments["namespace"].(string)
			result, err := s.cluster.GetContainersInPod(prowurl, podName, namespace)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_container_logs",
			mcp.WithDescription("Get logs of a specific container in a pod. Analyze these logs and print a succinct summary of important events, failures and errors if any."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("podName", mcp.Description("Pod name to fetch logs from"), mcp.Required()),
			mcp.WithString("namespace", mcp.Description("Namespace of the pod"), mcp.Required()),
			mcp.WithString("containerName", mcp.Description("Container name to fetch logs from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			podName := ctr.Params.Arguments["podName"].(string)
			namespace := ctr.Params.Arguments["namespace"].(string)
			containerName := ctr.Params.Arguments["containerName"].(string)
			result, err := s.cluster.GetContainerLogs(prowurl, podName, namespace, containerName)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_nodes_info",
			mcp.WithDescription("Get information of all nodes in the cluster. Format the output neatly with node name, architecture, OS image, kernel version, and other relevant details."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetNodesInfo(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_node_info_by_name",
			mcp.WithDescription("Get information of a specific node by name. Format the output neatly with node name, architecture, OS image, kernel version, and other relevant details."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("nodeName", mcp.Description("Node name to fetch information from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			nodeName := ctr.Params.Arguments["nodeName"].(string)
			result, err := s.cluster.GetNodeInfoByName(prowurl, nodeName)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_node_labels_by_name",
			mcp.WithDescription("Get labels of a specific node by name. Format the output neatly with node name and its labels."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("nodeName", mcp.Description("Node name to fetch labels from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			nodeName := ctr.Params.Arguments["nodeName"].(string)
			result, err := s.cluster.GetNodeLabelsByName(prowurl, nodeName)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_node_annotations_by_name",
			mcp.WithDescription("Get annotations of a specific node by name. Format the output neatly with node name and its annotations."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
			mcp.WithString("nodeName", mcp.Description("Node name to fetch annotations from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			nodeName := ctr.Params.Arguments["nodeName"].(string)
			result, err := s.cluster.GetNodeAnnotationsByName(prowurl, nodeName)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_nodes_labels",
			mcp.WithDescription("Get all labels from all nodes in the cluster. Format the output neatly with node name and its labels."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetNodesLabels(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_nodes_annotations",
			mcp.WithDescription("Get all annotations from all nodes in the cluster. Format the output neatly with node name and its annotations."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetNodesAnnotations(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_nodes_conditions",
			mcp.WithDescription("Get all conditions from all nodes in the cluster. Format the output neatly with node name and its conditions."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetNodesConditions(prowurl)
			return NewTextResult(result, err), nil
		}},
	}
}
