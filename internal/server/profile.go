package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerProfileTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_profile",
		Description: "Fetch the current Lexware profile for the configured API token.",
	}, s.getProfile)
}

func (s *Server) getProfile(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.GetProfile(ctx)
	return s.workflowResult("get profile", result, err)
}
