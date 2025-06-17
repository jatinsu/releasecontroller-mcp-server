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
			mcp.WithDescription("Get status summary of cluster operators. Clearly list the available, progressing, and degraded states of each operator."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetClusterOperatorStatusSummary(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("get_cluster_version_summary",
			mcp.WithDescription("Get the cluster version summary including the current version, desired version, and available updates."),
			mcp.WithString("prowurl", mcp.Description("Prow URL to fetch cluster version from"), mcp.Required()),
		), func(ctx context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.cluster.GetClusterVersionSummary(prowurl)
			return NewTextResult(result, err), nil
		}},
	}
}
