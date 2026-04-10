package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerCountryTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_countries",
		Description: "List all countries with their tax classifications (de, intraCommunity, thirdPartyCountry).",
	}, s.listCountries)
}

func (s *Server) listCountries(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.ListCountries(ctx)
	return s.workflowResult("list countries", result, err)
}
