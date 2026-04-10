package lexware

import "context"

func (c *Client) GetRecurringTemplate(ctx context.Context, id string) (*RecurringTemplateDetail, error) {
	result := &RecurringTemplateDetail{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/recurring-templates/" + id)
	if apiErr := wrapAPIError("get recurring template", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
