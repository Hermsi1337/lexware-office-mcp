package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type getDownPaymentInvoiceInput struct {
	ID string `json:"id" jsonschema:"UUID of the down payment invoice to retrieve"`
}

func (s *Server) registerDownPaymentInvoiceTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_down_payment_invoice",
		Description: "Retrieve a single down payment invoice by its UUID.",
	}, s.getDownPaymentInvoice)
}

func (s *Server) getDownPaymentInvoice(ctx context.Context, _ *mcp.CallToolRequest, input getDownPaymentInvoiceInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetDownPaymentInvoice(ctx, input.ID)
	return s.workflowResult("get down payment invoice", result, err)
}
