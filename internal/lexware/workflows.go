package lexware

import (
	"context"
	"strconv"
	"strings"
)

type CreateContactResult struct {
	ID string `json:"id"`
}

type CreateInvoiceResult struct {
	ID string `json:"id"`
}

type CreateQuotationResult struct {
	ID string `json:"id"`
}

type CreateCreditNoteResult struct {
	ID string `json:"id"`
}

type CreateArticleResult struct {
	ID string `json:"id"`
}

// ---------- Profile ----------

func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	result := &Profile{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/profile")
	if apiErr := wrapAPIError("get profile", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Contacts ----------

func (c *Client) CreateSimpleContact(ctx context.Context, name, sourceReference string) (*CreateContactResult, error) {
	contact := Contact{
		Version: 0,
		Roles: Roles{
			Customer: map[string]any{},
		},
		Person: Person{
			LastName: name,
		},
	}
	if strings.TrimSpace(sourceReference) != "" {
		contact.Note = sourceReference
	}

	result := &CreateContactResult{}
	resp, err := c.newRequest(ctx).
		SetBody(contact).
		SetResult(result).
		Post("/v1/contacts")
	if apiErr := wrapAPIError("create contact", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetContact(ctx context.Context, id string) (*ContactDetail, error) {
	result := &ContactDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/contacts/" + id)
	if apiErr := wrapAPIError("get contact", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ContactFilter holds optional query parameters for listing contacts.
type ContactFilter struct {
	Email    string
	Name     string
	Number   *int
	Customer *bool
	Vendor   *bool
	Page     int
}

func (c *Client) ListContacts(ctx context.Context, filter ContactFilter) (*Page[ContactDetail], error) {
	req := c.newRequest(ctx)
	if filter.Email != "" {
		req.SetQueryParam("email", filter.Email)
	}
	if filter.Name != "" {
		req.SetQueryParam("name", filter.Name)
	}
	if filter.Number != nil {
		req.SetQueryParam("number", strconv.Itoa(*filter.Number))
	}
	if filter.Customer != nil {
		req.SetQueryParam("customer", strconv.FormatBool(*filter.Customer))
	}
	if filter.Vendor != nil {
		req.SetQueryParam("vendor", strconv.FormatBool(*filter.Vendor))
	}
	if filter.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(filter.Page))
	}

	result := &Page[ContactDetail]{}
	resp, err := req.SetResult(result).Get("/v1/contacts")
	if apiErr := wrapAPIError("list contacts", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Invoices ----------

func (c *Client) CreateInvoice(ctx context.Context, invoice Invoice, finalize *bool) (*CreateInvoiceResult, error) {
	result := &CreateInvoiceResult{}
	resp, err := c.newRequest(ctx).
		SetBody(invoice).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/invoices")
	if apiErr := wrapAPIError("create invoice", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetInvoice(ctx context.Context, id string) (*InvoiceDetail, error) {
	result := &InvoiceDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/invoices/" + id)
	if apiErr := wrapAPIError("get invoice", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Articles ----------

func (c *Client) CreateArticle(ctx context.Context, article Article) (*CreateArticleResult, error) {
	result := &CreateArticleResult{}
	resp, err := c.newRequest(ctx).
		SetBody(article).
		SetResult(result).
		Post("/v1/articles")
	if apiErr := wrapAPIError("create article", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetArticle(ctx context.Context, id string) (*ArticleDetail, error) {
	result := &ArticleDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/articles/" + id)
	if apiErr := wrapAPIError("get article", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ArticleFilter holds optional query parameters for listing articles.
type ArticleFilter struct {
	ArticleNumber string
	Gtin          string
	Type          string
	Page          int
}

func (c *Client) ListArticles(ctx context.Context, filter ArticleFilter) (*Page[ArticleDetail], error) {
	req := c.newRequest(ctx)
	if filter.ArticleNumber != "" {
		req.SetQueryParam("articleNumber", filter.ArticleNumber)
	}
	if filter.Gtin != "" {
		req.SetQueryParam("gtin", filter.Gtin)
	}
	if filter.Type != "" {
		req.SetQueryParam("type", filter.Type)
	}
	if filter.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(filter.Page))
	}

	result := &Page[ArticleDetail]{}
	resp, err := req.SetResult(result).Get("/v1/articles")
	if apiErr := wrapAPIError("list articles", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Quotations ----------

func (c *Client) CreateQuotation(ctx context.Context, quotation Quotation, finalize *bool) (*CreateQuotationResult, error) {
	result := &CreateQuotationResult{}
	resp, err := c.newRequest(ctx).
		SetBody(quotation).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/quotations")
	if apiErr := wrapAPIError("create quotation", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetQuotation(ctx context.Context, id string) (*QuotationDetail, error) {
	result := &QuotationDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/quotations/" + id)
	if apiErr := wrapAPIError("get quotation", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Credit Notes ----------

func (c *Client) CreateCreditNote(ctx context.Context, creditNote CreditNote, finalize *bool) (*CreateCreditNoteResult, error) {
	result := &CreateCreditNoteResult{}
	resp, err := c.newRequest(ctx).
		SetBody(creditNote).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/credit-notes")
	if apiErr := wrapAPIError("create credit note", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetCreditNote(ctx context.Context, id string) (*CreditNoteDetail, error) {
	result := &CreditNoteDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/credit-notes/" + id)
	if apiErr := wrapAPIError("get credit note", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Countries ----------

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var result []Country
	resp, err := c.newRequest(ctx).
		SetResult(&result).
		Get("/v1/countries")
	if apiErr := wrapAPIError("list countries", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}

// ---------- Helpers ----------

func (c *Client) resolveFinalize(value *bool) bool {
	if value != nil {
		return *value
	}
	return c.finalizeInvoices
}
