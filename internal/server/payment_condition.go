package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) registerPaymentConditionTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_payment_conditions",
		Description: "List all payment condition presets configured in Lexware Office (e.g. Net 30, Immediate).",
	}, s.listPaymentConditions)
}

func (s *Server) listPaymentConditions(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.ListPaymentConditions(ctx)
	return s.workflowResult("list payment conditions", result, err)
}
