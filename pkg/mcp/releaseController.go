package mcp

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register the CLI tools for the release controller.
func (s *Server) initReleaseController() []server.ServerTool {
	return []server.ServerTool{
		{mcp.NewTool("list_release_controllers",
			mcp.WithDescription("Lists the available release controllers to use. Only two are available - OKD and OpenShift."),
		), s.listReleaseControllers},
		{mcp.NewTool("get_okd_release_controller",
			mcp.WithDescription("Gets the OKD release controller URL."),
		), func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return NewTextResult(s.releaseController.GetOKDReleaseController(), nil), nil
		}},
		{mcp.NewTool("get_ocp_release_controller",
			mcp.WithDescription("Gets the OpenShift release controller URL."),
		), func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return NewTextResult(s.releaseController.GetOCPReleaseController(), nil), nil
		}},
		{mcp.NewTool("list_release_streams",
			mcp.WithDescription("Lists all the release streams in the release controller."),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
		), s.listReleaseStreams},
		{mcp.NewTool("latest_accepted_release",
			mcp.WithDescription("Gets the latest accepted release for a given release stream."),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
		), s.latestAcceptedRelease},
		{mcp.NewTool("latest_rejected_release",
			mcp.WithDescription("Gets the latest rejected release for a given release stream."),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
		), s.latestRejectedRelease},
		{mcp.NewTool("list_failed_jobs_in_release",
			mcp.WithDescription("Lists all the failed jobs in a given release along with the prow job URL."),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
			mcp.WithString("tag", mcp.Description("The release tag"), mcp.Required()),
		), s.listFailedJobsInRelease},
		{mcp.NewTool("list_components_in_release",
			mcp.WithDescription("Lists the kubectl, kubernetes, coreos and tests versions in the release."),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
			mcp.WithString("tag", mcp.Description("The release tag"), mcp.Required()),
		), s.listComponentsInRelease},
		{mcp.NewTool("list_test_failures_for_release",
			mcp.WithDescription("Gets the failing tests for the particular job. List the failing tests in the release if there are any."),
			mcp.WithString("prowurl", mcp.Description("The prow job URL"), mcp.Required()),
		), func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prowurl := ctr.Params.Arguments["prowurl"].(string)
			result, err := s.releaseController.ListTestFailuresForRelease(prowurl)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("analyze_job_failures_for_release",
			mcp.WithDescription("Gets the build log file for the particular job. Analyze the job information and find out the reason for failure. Look for error strings in the log file. Print a short summary with relevant errors and keep it under 500 words"),
			mcp.WithString("prowurl", mcp.Description("The prow job URL"), mcp.Required()),
		), s.analyzeJobFailuresForRelease},
		{mcp.NewTool("list_features_from_updated_images_commits",
			mcp.WithDescription("Lists issues which are features from updated images commits - excludes OCPBUGS/CVEs"),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
			mcp.WithString("tag", mcp.Description("The release tag"), mcp.Required()),
		), func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
			stream := ctr.Params.Arguments["stream"].(string)
			tag := ctr.Params.Arguments["tag"].(string)
			result, err := s.releaseController.ListFeaturesFromUpdatedImagesCommits(releasecontroller, stream, tag)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("list_bugs_from_updated_images_commits",
			mcp.WithDescription("Lists issues which are bugs from updated images commits"),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
			mcp.WithString("tag", mcp.Description("The release tag"), mcp.Required()),
		), func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
			stream := ctr.Params.Arguments["stream"].(string)
			tag := ctr.Params.Arguments["tag"].(string)
			result, err := s.releaseController.ListBugsFromUpdatedImagesCommits(releasecontroller, stream, tag)
			return NewTextResult(result, err), nil
		}},
		{mcp.NewTool("list_cves_from_updated_images_commits",
			mcp.WithDescription("Lists issues which are CVEs from updated images commits"),
			mcp.WithString("releasecontroller", mcp.Description("The release controller host to query"), mcp.Required()),
			mcp.WithString("stream", mcp.Description("The release stream name"), mcp.Required()),
			mcp.WithString("tag", mcp.Description("The release tag"), mcp.Required()),
		), func(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
			stream := ctr.Params.Arguments["stream"].(string)
			tag := ctr.Params.Arguments["tag"].(string)
			result, err := s.releaseController.ListCVEsFromUpdatedImagesCommits(releasecontroller, stream, tag)
			return NewTextResult(result, err), nil
		}},
	}
}

func (s *Server) listReleaseControllers(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return NewTextResult(s.releaseController.ListReleaseControllers(), nil), nil
}

func (s *Server) listReleaseStreams(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
	result, err := s.releaseController.ListReleaseStreams(releasecontroller)
	return NewTextResult(result, err), nil
}

func (s *Server) latestAcceptedRelease(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
	stream := ctr.Params.Arguments["stream"].(string)
	result, err := s.releaseController.LatestAcceptedRelease(releasecontroller, stream)
	return NewTextResult(result, err), nil
}

func (s *Server) latestRejectedRelease(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
	stream := ctr.Params.Arguments["stream"].(string)
	result, err := s.releaseController.LatestRejectedRelease(releasecontroller, stream)
	return NewTextResult(result, err), nil
}

func (s *Server) listFailedJobsInRelease(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
	stream := ctr.Params.Arguments["stream"].(string)
	tag := ctr.Params.Arguments["tag"].(string)
	result, err := s.releaseController.ListFailedJobsInRelease(releasecontroller, stream, tag)
	return NewTextResult(result, err), nil
}

func (s *Server) listComponentsInRelease(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	releasecontroller := ctr.Params.Arguments["releasecontroller"].(string)
	stream := ctr.Params.Arguments["stream"].(string)
	tag := ctr.Params.Arguments["tag"].(string)
	result, err := s.releaseController.ListComponentsInRelease(releasecontroller, stream, tag)
	return NewTextResult(result, err), nil
}

func (s *Server) analyzeJobFailuresForRelease(_ context.Context, ctr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	prowurl := ctr.Params.Arguments["prowurl"].(string)
	result, err := s.releaseController.AnalyzeJobFailuresForRelease(prowurl)
	return NewTextResult(result, err), nil
}
