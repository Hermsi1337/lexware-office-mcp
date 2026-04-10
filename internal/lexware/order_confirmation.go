package lexware

import (
	"context"
	"strconv"
)

type CreateOrderConfirmationResult struct {
	ID string `json:"id"`
}

func (c *Client) CreateOrderConfirmation(ctx context.Context, oc OrderConfirmation, finalize *bool) (*CreateOrderConfirmationResult, error) {
	result := &CreateOrderConfirmationResult{}
	resp, err := c.newRequest(ctx).
		SetBody(oc).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/order-confirmations")
	if apiErr := wrapAPIError("create order confirmation", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetOrderConfirmation(ctx context.Context, id string) (*OrderConfirmationDetail, error) {
	result := &OrderConfirmationDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/order-confirmations/" + id)
	if apiErr := wrapAPIError("get order confirmation", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
