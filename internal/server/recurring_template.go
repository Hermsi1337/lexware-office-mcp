package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type getRecurringTemplateInput struct {
	ID string `json:"id" jsonschema:"UUID of the recurring template to retrieve"`
}

func (s *Server) registerRecurringTemplateTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_recurring_template",
		Description: "Retrieve a single recurring invoice template by its UUID, including schedule and line items.",
	}, s.getRecurringTemplate)
}

func (s *Server) getRecurringTemplate(ctx context.Context, _ *mcp.CallToolRequest, input getRecurringTemplateInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetRecurringTemplate(ctx, input.ID)
	return s.workflowResult("get recurring template", result, err)
}
