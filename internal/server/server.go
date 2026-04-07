package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type Server struct {
	*mcp.Server
	client *lexware.Client
}

type entityIDInput struct {
	ID string `json:"id" jsonschema:"UUID of the Lexware resource"`
}

type rawRequestInput struct {
	Method string            `json:"method,omitempty" jsonschema:"HTTP method, for example GET, POST, PUT or DELETE"`
	Path   string            `json:"path" jsonschema:"Lexware API path, for example /v1/contacts or /v1/invoices/{id}"`
	Query  map[string]string `json:"query,omitempty" jsonschema:"Optional query parameters"`
	Body   any               `json:"body,omitempty" jsonschema:"Optional JSON body for POST or PUT requests"`
	Accept string            `json:"accept,omitempty" jsonschema:"Optional Accept header. Defaults to application/json"`
}

type pageInput struct {
	Page int `json:"page,omitempty" jsonschema:"Page number starting at 0"`
}

type createEntityInput struct {
	Payload map[string]any `json:"payload" jsonschema:"JSON payload exactly as required by the Lexware API"`
}

type updateEntityInput struct {
	ID      string         `json:"id" jsonschema:"UUID of the Lexware resource"`
	Payload map[string]any `json:"payload" jsonschema:"Full JSON payload including version when required by Lexware"`
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
		Name:        "lexware_list_contacts",
		Description: "List contacts from Lexware with optional paging.",
	}, s.listContacts)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_contact",
		Description: "Fetch a single contact by id.",
	}, s.getContact)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_contact",
		Description: "Create a contact with a raw Lexware JSON payload.",
	}, s.createContact)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_update_contact",
		Description: "Update a contact by id with a raw Lexware JSON payload.",
	}, s.updateContact)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_articles",
		Description: "List articles from Lexware with optional paging.",
	}, s.listArticles)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_article",
		Description: "Fetch a single article by id.",
	}, s.getArticle)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_invoice",
		Description: "Fetch a single invoice by id.",
	}, s.getInvoice)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_vouchers",
		Description: "List vouchers with optional paging.",
	}, s.listVouchers)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_api_request",
		Description: "Make a raw authenticated request against the Lexware API for endpoints not yet wrapped by dedicated tools.",
	}, s.rawRequest)
}

func (s *Server) getProfile(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/profile"})
}

func (s *Server) listContacts(ctx context.Context, _ *mcp.CallToolRequest, input pageInput) (*mcp.CallToolResult, map[string]any, error) {
	query := pagingQuery(input.Page)
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/contacts", Query: query})
}

func (s *Server) getContact(ctx context.Context, _ *mcp.CallToolRequest, input entityIDInput) (*mcp.CallToolResult, map[string]any, error) {
	id, err := requireID(input.ID)
	if err != nil {
		return nil, nil, err
	}
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/contacts/" + id})
}

func (s *Server) createContact(ctx context.Context, _ *mcp.CallToolRequest, input createEntityInput) (*mcp.CallToolResult, map[string]any, error) {
	return s.do(ctx, lexware.Request{Method: "POST", Path: "/v1/contacts", Body: input.Payload})
}

func (s *Server) updateContact(ctx context.Context, _ *mcp.CallToolRequest, input updateEntityInput) (*mcp.CallToolResult, map[string]any, error) {
	id, err := requireID(input.ID)
	if err != nil {
		return nil, nil, err
	}
	return s.do(ctx, lexware.Request{Method: "PUT", Path: "/v1/contacts/" + id, Body: input.Payload})
}

func (s *Server) listArticles(ctx context.Context, _ *mcp.CallToolRequest, input pageInput) (*mcp.CallToolResult, map[string]any, error) {
	query := pagingQuery(input.Page)
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/articles", Query: query})
}

func (s *Server) getArticle(ctx context.Context, _ *mcp.CallToolRequest, input entityIDInput) (*mcp.CallToolResult, map[string]any, error) {
	id, err := requireID(input.ID)
	if err != nil {
		return nil, nil, err
	}
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/articles/" + id})
}

func (s *Server) getInvoice(ctx context.Context, _ *mcp.CallToolRequest, input entityIDInput) (*mcp.CallToolResult, map[string]any, error) {
	id, err := requireID(input.ID)
	if err != nil {
		return nil, nil, err
	}
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/invoices/" + id})
}

func (s *Server) listVouchers(ctx context.Context, _ *mcp.CallToolRequest, input pageInput) (*mcp.CallToolResult, map[string]any, error) {
	query := pagingQuery(input.Page)
	return s.do(ctx, lexware.Request{Method: "GET", Path: "/v1/voucherlist", Query: query})
}

func (s *Server) rawRequest(ctx context.Context, _ *mcp.CallToolRequest, input rawRequestInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.Path) == "" {
		return nil, nil, fmt.Errorf("path is required")
	}

	return s.do(ctx, lexware.Request{
		Method: input.Method,
		Path:   input.Path,
		Query:  input.Query,
		Body:   input.Body,
		Accept: input.Accept,
	})
}

func (s *Server) do(ctx context.Context, req lexware.Request) (*mcp.CallToolResult, map[string]any, error) {
	resp, err := s.client.Do(ctx, req)
	if err != nil && resp == nil {
		return nil, nil, err
	}

	out := map[string]any{
		"statusCode":  resp.StatusCode,
		"contentType": resp.ContentType,
		"headers":     resp.Headers,
	}
	if resp.Body != nil {
		out["body"] = resp.Body
	}
	if resp.BodyText != "" {
		out["bodyText"] = resp.BodyText
	}
	if resp.BodyBase64 != "" {
		out["bodyBase64"] = resp.BodyBase64
	}

	if err != nil {
		text := fmt.Sprintf("Lexware API request failed with status %d", resp.StatusCode)
		if b, marshalErr := json.MarshalIndent(out, "", "  "); marshalErr == nil {
			text += "\n" + string(b)
		}
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: text},
			},
		}, out, nil
	}

	return nil, out, nil
}

func pagingQuery(page int) map[string]string {
	if page <= 0 {
		return nil
	}
	return map[string]string{"page": fmt.Sprintf("%d", page)}
}

func requireID(id string) (string, error) {
	value := strings.TrimSpace(id)
	if value == "" {
		return "", fmt.Errorf("id is required")
	}
	return value, nil
}
