package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerPostingCategoryTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_posting_categories",
		Description: "List all posting categories used for voucher bookkeeping in Lexware Office.",
	}, s.listPostingCategories)
}

func (s *Server) listPostingCategories(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.ListPostingCategories(ctx)
	return s.workflowResult("list posting categories", result, err)
}
