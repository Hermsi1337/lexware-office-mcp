package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type createSimpleContactInput struct {
	Name            string `json:"name" jsonschema:"Display name used as the contact last name in the simple contact helper"`
	SourceReference string `json:"sourceReference,omitempty" jsonschema:"Optional source reference stored in the contact note, for example an order id"`
}

type getContactInput struct {
	ID string `json:"id" jsonschema:"UUID of the contact to retrieve"`
}

type listContactsInput struct {
	Email    string `json:"email,omitempty" jsonschema:"Filter contacts by email (min 3 chars, substring match)"`
	Name     string `json:"name,omitempty" jsonschema:"Filter contacts by name (min 3 chars, supports % and _ wildcards)"`
	Number   *int   `json:"number,omitempty" jsonschema:"Filter contacts by customer or vendor number"`
	Customer *bool  `json:"customer,omitempty" jsonschema:"Filter for contacts with customer role"`
	Vendor   *bool  `json:"vendor,omitempty" jsonschema:"Filter for contacts with vendor role"`
	Page     int    `json:"page,omitempty" jsonschema:"Page number for pagination (0-based)"`
}

func (s *Server) registerContactTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_simple_contact",
		Description: "Create a simple customer contact with a last name and optional source reference.",
	}, s.createSimpleContact)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_contact",
		Description: "Retrieve a single contact by its UUID.",
	}, s.getContact)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_contacts",
		Description: "List contacts with optional filters and pagination. Supports name search with wildcards (% and _, min 3 chars), email substring match (min 3 chars), exact customer/vendor number, and role filtering. This is the closest thing to a contact search the Lexware API offers.",
	}, s.listContacts)
}

func (s *Server) createSimpleContact(ctx context.Context, _ *mcp.CallToolRequest, input createSimpleContactInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, nil, fmt.Errorf("name is required")
	}

	result, err := s.client.CreateSimpleContact(ctx, input.Name, input.SourceReference)
	return s.workflowResult("create contact", result, err)
}

func (s *Server) getContact(ctx context.Context, _ *mcp.CallToolRequest, input getContactInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetContact(ctx, input.ID)
	return s.workflowResult("get contact", result, err)
}

func (s *Server) listContacts(ctx context.Context, _ *mcp.CallToolRequest, input listContactsInput) (*mcp.CallToolResult, map[string]any, error) {
	filter := lexware.ContactFilter{
		Email:    input.Email,
		Name:     input.Name,
		Number:   input.Number,
		Customer: input.Customer,
		Vendor:   input.Vendor,
		Page:     input.Page,
	}

	result, err := s.client.ListContacts(ctx, filter)
	return s.workflowResult("list contacts", result, err)
}
