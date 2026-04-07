package lexware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CreateContactResult struct {
	ID string `json:"id"`
}

type CreateInvoiceResult struct {
	ID string `json:"id"`
}

func (c *Client) GetProfile(ctx context.Context) (*Profile, *Response, error) {
	result := &Profile{}
	resp, err := c.restClient.R().
		SetContext(ctx).
		SetResult(result).
		Get("/v1/profile")
	if err != nil {
		return nil, nil, fmt.Errorf("send request: %w", err)
	}

	rawResp, decodeErr := decodeResponse(resp.RawResponse, resp.Body())
	if decodeErr != nil {
		return nil, nil, decodeErr
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return nil, rawResp, fmt.Errorf("lexware api returned status %d", resp.StatusCode())
	}

	return result, rawResp, nil
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

	result := &CreateContactResult{}
	resp, err := c.restClient.R().
		SetContext(ctx).
		SetBody(contact).
		SetResult(result).
		Post("/v1/contacts")
	if err != nil {
		return nil, nil, fmt.Errorf("send request: %w", err)
	}

	rawResp, decodeErr := decodeResponse(resp.RawResponse, resp.Body())
	if decodeErr != nil {
		return nil, nil, decodeErr
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return result, rawResp, fmt.Errorf("lexware api returned status %d", resp.StatusCode())
	}

	return result, rawResp, nil
}

func (c *Client) CreateInvoice(ctx context.Context, invoice Invoice, finalize *bool) (*CreateInvoiceResult, *Response, error) {
	result := &CreateInvoiceResult{}
	resp, err := c.restClient.R().
		SetContext(ctx).
		SetBody(invoice).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/invoices")
	if err != nil {
		return nil, nil, err
	}

	rawResp, decodeErr := decodeResponse(resp.RawResponse, resp.Body())
	if decodeErr != nil {
		return nil, nil, decodeErr
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return result, rawResp, fmt.Errorf("lexware api returned status %d", resp.StatusCode())
	}

	return result, rawResp, nil
}

func (c *Client) resolveFinalize(value *bool) bool {
	if value != nil {
		return *value
	}
	return c.finalizeInvoices
}
