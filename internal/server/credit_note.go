package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createCreditNoteInput struct {
	CreditNote lexware.CreditNote `json:"creditNote" jsonschema:"Credit note payload"`
	Finalize   *bool              `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getCreditNoteInput struct {
	ID string `json:"id" jsonschema:"UUID of the credit note to retrieve"`
}

func (s *Server) registerCreditNoteTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_credit_note",
		Description: "Create a credit note with line items, address, and an optional finalize flag.",
	}, s.createCreditNote)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_credit_note",
		Description: "Retrieve a single credit note by its UUID.",
	}, s.getCreditNote)
}

func (s *Server) createCreditNote(ctx context.Context, _ *mcp.CallToolRequest, input createCreditNoteInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateCreditNote(ctx, input.CreditNote, input.Finalize)
	return s.workflowResult("create credit note", result, err)
}

func (s *Server) getCreditNote(ctx context.Context, _ *mcp.CallToolRequest, input getCreditNoteInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetCreditNote(ctx, input.ID)
	return s.workflowResult("get credit note", result, err)
}
