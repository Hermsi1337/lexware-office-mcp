package lexware

import (
	"context"
	"strconv"
)

type CreateDeliveryNoteResult struct {
	ID string `json:"id"`
}

func (c *Client) CreateDeliveryNote(ctx context.Context, note DeliveryNote, finalize *bool) (*CreateDeliveryNoteResult, error) {
	result := &CreateDeliveryNoteResult{}
	resp, err := c.newRequest(ctx).
		SetBody(note).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/delivery-notes")
	if apiErr := wrapAPIError("create delivery note", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetDeliveryNote(ctx context.Context, id string) (*DeliveryNoteDetail, error) {
	result := &DeliveryNoteDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/delivery-notes/" + id)
	if apiErr := wrapAPIError("get delivery note", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
