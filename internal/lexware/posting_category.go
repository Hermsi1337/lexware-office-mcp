package lexware

import "context"

func (c *Client) ListPostingCategories(ctx context.Context) ([]PostingCategory, error) {
	var result []PostingCategory
	resp, err := c.newRequest(ctx).
		SetResult(&result).
		Get("/v1/posting-categories")
	if apiErr := wrapAPIError("list posting categories", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
