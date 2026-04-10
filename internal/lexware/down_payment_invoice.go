package lexware

import "context"

func (c *Client) GetDownPaymentInvoice(ctx context.Context, id string) (*DownPaymentInvoiceDetail, error) {
	result := &DownPaymentInvoiceDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/down-payment-invoices/" + id)
	if apiErr := wrapAPIError("get down payment invoice", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
