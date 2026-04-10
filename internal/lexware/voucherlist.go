package lexware

import (
	"context"
	"strconv"
)

// VoucherlistFilter holds query parameters for the voucherlist endpoint.
type VoucherlistFilter struct {
	VoucherType   string // Comma-separated: salesinvoice,salescreditnote,purchaseinvoice,purchasecreditnote
	VoucherStatus string // Comma-separated: open,paid,paidoff,voided,transferred,sepadebit,unchecked
	Page          int
	Size          int // 1-250, default 250
}

func (c *Client) ListVouchers(ctx context.Context, filter VoucherlistFilter) (*Page[VoucherListItem], error) {
	req := c.newRequest(ctx)
	if filter.VoucherType != "" {
		req.SetQueryParam("voucherType", filter.VoucherType)
	}
	if filter.VoucherStatus != "" {
		req.SetQueryParam("voucherStatus", filter.VoucherStatus)
	}
	if filter.Page > 0 {
		req.SetQueryParam("page", strconv.Itoa(filter.Page))
	}
	if filter.Size > 0 {
		req.SetQueryParam("size", strconv.Itoa(filter.Size))
	}

	result := &Page[VoucherListItem]{}
	resp, err := req.SetResult(result).Get("/v1/voucherlist")
	if apiErr := wrapAPIError("list vouchers", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
