package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createDeliveryNoteInput struct {
	DeliveryNote lexware.DeliveryNote `json:"deliveryNote" jsonschema:"Delivery note payload"`
	Finalize     *bool                `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getDeliveryNoteInput struct {
	ID string `json:"id" jsonschema:"UUID of the delivery note to retrieve"`
}

func (s *Server) registerDeliveryNoteTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_delivery_note",
		Description: "Create a delivery note with line items and an optional finalize flag.",
	}, s.createDeliveryNote)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_delivery_note",
		Description: "Retrieve a single delivery note by its UUID.",
	}, s.getDeliveryNote)
}

func (s *Server) createDeliveryNote(ctx context.Context, _ *mcp.CallToolRequest, input createDeliveryNoteInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateDeliveryNote(ctx, input.DeliveryNote, input.Finalize)
	return s.workflowResult("create delivery note", result, err)
}

func (s *Server) getDeliveryNote(ctx context.Context, _ *mcp.CallToolRequest, input getDeliveryNoteInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetDeliveryNote(ctx, input.ID)
	return s.workflowResult("get delivery note", result, err)
}
