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

func (c *Client) resolveFinalize(value *bool) bool {
	if value != nil {
		return *value
	}
	return c.finalizeInvoices
}
