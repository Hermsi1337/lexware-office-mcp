package lexware

import (
	"context"
	"strconv"
)

type CreateQuotationResult struct {
	ID string `json:"id"`
}

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
