package server

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
	"github.com/dennis/lexware-office-mcp/internal/version"
)

type Server struct {
	*mcp.Server
	client *lexware.Client
}

func New(client *lexware.Client) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "lexware-office-mcp",
		Version: version.Version,
	}, nil)

	wrapped := &Server{
		Server: srv,
		client: client,
	}

	wrapped.registerTools()
	return srv
}

func (s *Server) registerTools() {
	s.registerProfileTools()
	s.registerContactTools()
	s.registerInvoiceTools()
	s.registerArticleTools()
	s.registerQuotationTools()
	s.registerCreditNoteTools()
	s.registerVoucherlistTools()
	s.registerDeliveryNoteTools()
	s.registerOrderConfirmationTools()
	s.registerCountryTools()
	s.registerPaymentConditionTools()
	s.registerPostingCategoryTools()
	s.registerDownPaymentInvoiceTools()
	s.registerRecurringTemplateTools()
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
