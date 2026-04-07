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

func (c *Client) CreateSimpleContact(ctx context.Context, name, sourceReference string) (*CreateContactResult, *Response, error) {
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

	resp, err := c.Do(ctx, Request{
		Method: "POST",
		Path:   "/v1/contacts",
		Body:   contact,
	})
	if err != nil && resp == nil {
		return nil, nil, err
	}

	result := &CreateContactResult{}
	if resp != nil && resp.Body != nil {
		if id, ok := extractID(resp.Body); ok {
			result.ID = id
		}
	}

	return result, resp, err
}

func (c *Client) CreateInvoice(ctx context.Context, invoice Invoice, finalize *bool) (*CreateInvoiceResult, *Response, error) {
	resp, err := c.Do(ctx, Request{
		Method: "POST",
		Path:   "/v1/invoices",
		Query: map[string]string{
			"finalize": strconv.FormatBool(c.resolveFinalize(finalize)),
		},
		Body: invoice,
	})
	if err != nil && resp == nil {
		return nil, nil, err
	}

	result := &CreateInvoiceResult{}
	if resp != nil && resp.Body != nil {
		if id, ok := extractID(resp.Body); ok {
			result.ID = id
		}
	}

	return result, resp, err
}

func (c *Client) resolveFinalize(value *bool) bool {
	if value != nil {
		return *value
	}
	return c.finalizeInvoices
}

func extractID(body any) (string, bool) {
	asMap, ok := body.(map[string]any)
	if !ok {
		return "", false
	}

	rawID, ok := asMap["id"]
	if !ok {
		return "", false
	}

	id, ok := rawID.(string)
	return id, ok && strings.TrimSpace(id) != ""
}
