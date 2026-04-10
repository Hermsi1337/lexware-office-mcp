package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
	"github.com/dennis/lexware-office-mcp/internal/version"
)

type Server struct {
	*mcp.Server
	client *lexware.Client
}

// ---------- Input types ----------

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

type createInvoiceInput struct {
	Invoice  lexware.Invoice `json:"invoice" jsonschema:"Invoice payload"`
	Finalize *bool           `json:"finalize,omitempty" jsonschema:"Optional override for Lexware invoice finalization"`
}

type getInvoiceInput struct {
	ID string `json:"id" jsonschema:"UUID of the invoice to retrieve"`
}

type createArticleInput struct {
	Article lexware.Article `json:"article" jsonschema:"Article payload with title, type (PRODUCT or SERVICE), unitName, and price"`
}

type getArticleInput struct {
	ID string `json:"id" jsonschema:"UUID of the article to retrieve"`
}

type listArticlesInput struct {
	ArticleNumber string `json:"articleNumber,omitempty" jsonschema:"Filter by exact article number"`
	Gtin          string `json:"gtin,omitempty" jsonschema:"Filter by GTIN"`
	Type          string `json:"type,omitempty" jsonschema:"Filter by type: PRODUCT or SERVICE"`
	Page          int    `json:"page,omitempty" jsonschema:"Page number for pagination (0-based)"`
}

type createQuotationInput struct {
	Quotation lexware.Quotation `json:"quotation" jsonschema:"Quotation payload"`
	Finalize  *bool             `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getQuotationInput struct {
	ID string `json:"id" jsonschema:"UUID of the quotation to retrieve"`
}

type createCreditNoteInput struct {
	CreditNote lexware.CreditNote `json:"creditNote" jsonschema:"Credit note payload"`
	Finalize   *bool              `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getCreditNoteInput struct {
	ID string `json:"id" jsonschema:"UUID of the credit note to retrieve"`
}

type listVouchersInput struct {
	VoucherType   string `json:"voucherType,omitempty" jsonschema:"Comma-separated voucher types: salesinvoice, salescreditnote, purchaseinvoice, purchasecreditnote"`
	VoucherStatus string `json:"voucherStatus,omitempty" jsonschema:"Comma-separated statuses: open, paid, paidoff, voided, transferred, sepadebit, unchecked"`
	Page          int    `json:"page,omitempty" jsonschema:"Page number for pagination (0-based)"`
	Size          int    `json:"size,omitempty" jsonschema:"Results per page (1-250, default 250)"`
}

type createDeliveryNoteInput struct {
	DeliveryNote lexware.DeliveryNote `json:"deliveryNote" jsonschema:"Delivery note payload"`
	Finalize     *bool                `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getDeliveryNoteInput struct {
	ID string `json:"id" jsonschema:"UUID of the delivery note to retrieve"`
}

type createOrderConfirmationInput struct {
	OrderConfirmation lexware.OrderConfirmation `json:"orderConfirmation" jsonschema:"Order confirmation payload"`
	Finalize          *bool                     `json:"finalize,omitempty" jsonschema:"Optional override for Lexware finalization"`
}

type getOrderConfirmationInput struct {
	ID string `json:"id" jsonschema:"UUID of the order confirmation to retrieve"`
}

// ---------- Constructor ----------

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

// ---------- Tool registration ----------

func (s *Server) registerTools() {
	// Profile
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_profile",
		Description: "Fetch the current Lexware profile for the configured API token.",
	}, s.getProfile)

	// Contacts
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

	// Invoices
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_invoice",
		Description: "Create an invoice with line items, address, tax conditions, and an optional finalize flag.",
	}, s.createInvoice)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_invoice",
		Description: "Retrieve a single invoice by its UUID, including status, line items, and totals.",
	}, s.getInvoice)

	// Articles
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_article",
		Description: "Create an article (product or service) with title, unit, and pricing.",
	}, s.createArticle)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_article",
		Description: "Retrieve a single article by its UUID.",
	}, s.getArticle)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_articles",
		Description: "List articles with optional filters and pagination. Supports exact match on article number and GTIN, and filtering by type (PRODUCT or SERVICE). No full-text search available.",
	}, s.listArticles)

	// Quotations
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_quotation",
		Description: "Create a quotation with line items, address, and an optional finalize flag.",
	}, s.createQuotation)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_quotation",
		Description: "Retrieve a single quotation by its UUID.",
	}, s.getQuotation)

	// Credit Notes
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_credit_note",
		Description: "Create a credit note with line items, address, and an optional finalize flag.",
	}, s.createCreditNote)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_credit_note",
		Description: "Retrieve a single credit note by its UUID.",
	}, s.getCreditNote)

	// Voucherlist
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_vouchers",
		Description: "List vouchers across all document types with filters for voucher type and status. Returns a unified view of invoices, credit notes, and purchase documents with contact references. Note: no search by recipient name is available; to find vouchers for a specific contact, first look up the contact via lexware_list_contacts.",
	}, s.listVouchers)

	// Delivery Notes
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_delivery_note",
		Description: "Create a delivery note with line items and an optional finalize flag.",
	}, s.createDeliveryNote)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_delivery_note",
		Description: "Retrieve a single delivery note by its UUID.",
	}, s.getDeliveryNote)

	// Order Confirmations
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_create_order_confirmation",
		Description: "Create an order confirmation with line items and an optional finalize flag.",
	}, s.createOrderConfirmation)

	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_get_order_confirmation",
		Description: "Retrieve a single order confirmation by its UUID.",
	}, s.getOrderConfirmation)

	// Countries
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_countries",
		Description: "List all countries with their tax classifications (de, intraCommunity, thirdPartyCountry).",
	}, s.listCountries)
}

// ---------- Handlers ----------

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

func (s *Server) createArticle(ctx context.Context, _ *mcp.CallToolRequest, input createArticleInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateArticle(ctx, input.Article)
	return s.workflowResult("create article", result, err)
}

func (s *Server) getArticle(ctx context.Context, _ *mcp.CallToolRequest, input getArticleInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetArticle(ctx, input.ID)
	return s.workflowResult("get article", result, err)
}

func (s *Server) listArticles(ctx context.Context, _ *mcp.CallToolRequest, input listArticlesInput) (*mcp.CallToolResult, map[string]any, error) {
	filter := lexware.ArticleFilter{
		ArticleNumber: input.ArticleNumber,
		Gtin:          input.Gtin,
		Type:          input.Type,
		Page:          input.Page,
	}

	result, err := s.client.ListArticles(ctx, filter)
	return s.workflowResult("list articles", result, err)
}

func (s *Server) createQuotation(ctx context.Context, _ *mcp.CallToolRequest, input createQuotationInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.CreateQuotation(ctx, input.Quotation, input.Finalize)
	return s.workflowResult("create quotation", result, err)
}

func (s *Server) getQuotation(ctx context.Context, _ *mcp.CallToolRequest, input getQuotationInput) (*mcp.CallToolResult, map[string]any, error) {
	if strings.TrimSpace(input.ID) == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	result, err := s.client.GetQuotation(ctx, input.ID)
	return s.workflowResult("get quotation", result, err)
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

func (s *Server) listVouchers(ctx context.Context, _ *mcp.CallToolRequest, input listVouchersInput) (*mcp.CallToolResult, map[string]any, error) {
	filter := lexware.VoucherlistFilter{
		VoucherType:   input.VoucherType,
		VoucherStatus: input.VoucherStatus,
		Page:          input.Page,
		Size:          input.Size,
	}

	result, err := s.client.ListVouchers(ctx, filter)
	return s.workflowResult("list vouchers", result, err)
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

func (s *Server) listCountries(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.client.ListCountries(ctx)
	return s.workflowResult("list countries", result, err)
}

// ---------- Result helper ----------

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
