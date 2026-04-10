package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createOrderConfirmationInput struct {
	OrderConfirmation lexware.OrderConfirmation `json:"orderConfirmation" jsonschema:"Order confirmation payload"`
	Finalize          *bool                     `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getOrderConfirmationInput struct {
	ID string `json:"id" jsonschema:"UUID of the order confirmation to retrieve"`
}

func (s *Server) registerOrderConfirmationTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_order_confirmation",
		Description: "Create an order confirmation with line items and an optional finalize flag.",
	}, s.createOrderConfirmation)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_order_confirmation",
		Description: "Retrieve a single order confirmation by its UUID.",
	}, s.getOrderConfirmation)
}

func (s *Server) createOrderConfirmation(ctx context.Context, _ *mcp.CallToolRequest, input createOrderConfirmationInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateOrderConfirmation(ctx, input.OrderConfirmation, input.Finalize)
	return s.workflowResult("create order confirmation", result, err)
}

func (s *Server) getOrderConfirmation(ctx context.Context, _ *mcp.CallToolRequest, input getOrderConfirmationInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetOrderConfirmation(ctx, input.ID)
	return s.workflowResult("get order confirmation", result, err)
}
