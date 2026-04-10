package lexware

import "context"

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var result []Country
	resp, err := c.newRequest(ctx).
		SetResult(&result).
		Get("/v1/countries")
	if apiErr := wrapAPIError("list countries", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
