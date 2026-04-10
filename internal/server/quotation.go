package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createQuotationInput struct {
	Quotation lexware.Quotation `json:"quotation" jsonschema:"Quotation payload"`
	Finalize  *bool             `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getQuotationInput struct {
	ID string `json:"id" jsonschema:"UUID of the quotation to retrieve"`
}

func (s *Server) registerQuotationTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_quotation",
		Description: "Create a quotation with line items, address, and an optional finalize flag.",
	}, s.createQuotation)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_quotation",
		Description: "Retrieve a single quotation by its UUID.",
	}, s.getQuotation)
}

func (s *Server) createQuotation(ctx context.Context, _ *mcp.CallToolRequest, input createQuotationInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateQuotation(ctx, input.Quotation, input.Finalize)
	return s.workflowResult("create quotation", result, err)
}

func (s *Server) getQuotation(ctx context.Context, _ *mcp.CallToolRequest, input getQuotationInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetQuotation(ctx, input.ID)
	return s.workflowResult("get quotation", result, err)
}
