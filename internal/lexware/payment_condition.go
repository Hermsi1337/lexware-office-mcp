package lexware

import "context"

func (c *Client) ListPaymentConditions(ctx context.Context) ([]PaymentConditionItem, error) {
	var result []PaymentConditionItem
	resp, err := c.newRequest(ctx).
		SetResult(&result).
		Get("/v1/payment-conditions")
	if apiErr := wrapAPIError("list payment conditions", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
