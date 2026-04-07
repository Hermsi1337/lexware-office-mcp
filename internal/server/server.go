package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type Server struct {
	*mcp.Server
	client *lexware.Client
}

type createSimpleContactInput struct {
	Name            string `json:"name" jsonschema:"Display name used as the contact last name in the simple contact helper"`
	SourceReference string `json:"sourceReference,omitempty" jsonschema:"Optional source reference stored in the contact note, for example an order id"`
}

type createInvoiceInput struct {
	Invoice  lexware.Invoice `json:"invoice" jsonschema:"Invoice payload based on the legacy project structure"`
	Finalize *bool           `json:"finalize,omitempty" jsonschema:"Optional override for Lexware invoice finalization"`
}

func New(client *lexware.Client) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "lexware-office-mcp",
		Version: "0.1.0",
	}, nil)

	wrapped := &Server{
		Server: srv,
		client: client,
	}

	wrapped.registerTools()
	return srv
}

func (s *Server) registerTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_profile",
		Description: "Fetch the current Lexware profile for the configured API token.",
	}, s.getProfile)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_simple_contact",
		Description: "Create a simple customer contact using the legacy integration shape.",
	}, s.createSimpleContact)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_invoice",
		Description: "Create an invoice using the legacy integration invoice shape and an optional finalize flag.",
	}, s.createInvoice)
}

func (s *Server) getProfile(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.GetProfile(ctx)
	return s.workflowResult("get profile", result, err)
}

func (s *Server) createSimpleContact(ctx context.Context, _ *mcp.CallToolRequest, input createSimpleContactInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, nil, fmt.Errorf("name is required")
	}

	result, err := s.client.CreateSimpleContact(ctx, input.Name, input.SourceReference)
	return s.workflowResult("create contact", result, err)
}

func (s *Server) createInvoice(ctx context.Context, _ *mcp.CallToolRequest, input createInvoiceInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateInvoice(ctx, input.Invoice, input.Finalize)
	return s.workflowResult("create invoice", result, err)
}

func (s *Server) workflowResult(action string, result any, err error) (*mcp.CallToolResult, map[string]any, error) {
	payload := map[string]any{
		"result": result,
	}

	if err != nil {
		text := fmt.Sprintf("Lexware %s failed", action)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: text + ": " + err.Error()},
			},
		}, payload, nil
	}

	return nil, payload, nil
}
