package lexware

import (
	"context"
	"strconv"
)

type CreateInvoiceResult struct {
	ID string `json:"id"`
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
