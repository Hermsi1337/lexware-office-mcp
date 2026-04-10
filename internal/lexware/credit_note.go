package lexware

import (
	"context"
	"strconv"
)

type CreateCreditNoteResult struct {
	ID string `json:"id"`
}

func (c *Client) CreateCreditNote(ctx context.Context, creditNote CreditNote, finalize *bool) (*CreateCreditNoteResult, error) {
	result := &CreateCreditNoteResult{}
	resp, err := c.newRequest(ctx).
		SetBody(creditNote).
		SetResult(result).
		SetQueryParam("finalize", strconv.FormatBool(c.resolveFinalize(finalize))).
		Post("/v1/credit-notes")
	if apiErr := wrapAPIError("create credit note", resp, err); apiErr != nil {
		return result, apiErr
	}

	return result, nil
}

func (c *Client) GetCreditNote(ctx context.Context, id string) (*CreditNoteDetail, error) {
	result := &CreditNoteDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/credit-notes/" + id)
	if apiErr := wrapAPIError("get credit note", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
