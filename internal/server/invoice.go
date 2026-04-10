package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createInvoiceInput struct {
	Invoice  lexware.Invoice `json:"invoice" jsonschema:"Invoice payload"`
	Finalize *bool           `json:"finalize,omitempty" jsonschema:"Optional override for Lexware invoice finalization"`
}

type getInvoiceInput struct {
	ID string `json:"id" jsonschema:"UUID of the invoice to retrieve"`
}

func (s *Server) registerInvoiceTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_invoice",
		Description: "Create an invoice with line items, address, tax conditions, and an optional finalize flag.",
	}, s.createInvoice)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_invoice",
		Description: "Retrieve a single invoice by its UUID, including status, line items, and totals.",
	}, s.getInvoice)
}

func (s *Server) createInvoice(ctx context.Context, _ *mcp.CallToolRequest, input createInvoiceInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateInvoice(ctx, input.Invoice, input.Finalize)
	return s.workflowResult("create invoice", result, err)
}

func (s *Server) getInvoice(ctx context.Context, _ *mcp.CallToolRequest, input getInvoiceInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetInvoice(ctx, input.ID)
	return s.workflowResult("get invoice", result, err)
}
