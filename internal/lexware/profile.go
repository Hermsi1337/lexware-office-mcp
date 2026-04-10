package lexware

import "context"

func (c *Client) GetProfile(ctx context.Context) (*Profile, error) {
	result := &Profile{}
	resp, err := c.newRequest(ctx).
		SetResult(result).
		Get("/v1/profile")
	if apiErr := wrapAPIError("get profile", resp, err); apiErr != nil {
		return nil, apiErr
	}

	return result, nil
}
