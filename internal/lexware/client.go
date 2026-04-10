package lexware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	restClient       *resty.Client
	finalizeInvoices bool
}

func NewClient(cfg Config) *Client {
	restClient := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.HTTPTimeout).
		SetHeader("Authorization", "Bearer "+cfg.APIToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", cfg.UserAgent).
		SetRetryCount(5).
		SetRetryWaitTime(10 * time.Second).
		SetRetryMaxWaitTime(60 * time.Second).
		AddRetryCondition(func(resp *resty.Response, err error) bool {
			return err == nil && resp != nil && resp.StatusCode() == http.StatusTooManyRequests
		})

	return &Client{
		restClient:       restClient,
		finalizeInvoices: cfg.FinalizeInvoices,
	}
}

func wrapAPIError(action string, resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("%s request failed: %w", action, err)
	}
	if resp == nil {
		return fmt.Errorf("%s failed without a response", action)
	}
	if resp.StatusCode() < http.StatusBadRequest {
		return nil
	}

	body := strings.TrimSpace(resp.String())
	if body == "" {
		return fmt.Errorf("%s failed with status %d", action, resp.StatusCode())
	}

	return fmt.Errorf("%s failed with status %d: %s", action, resp.StatusCode(), body)
}

func (c *Client) FinalizeInvoices() bool {
	return c.finalizeInvoices
}

func (c *Client) newRequest(ctx context.Context) *resty.Request {
	return c.restClient.R().SetContext(ctx)
}

func (c *Client) resolveFinalize(value *bool) bool {
	if value != nil {
		return *value
	}
	return c.finalizeInvoices
}
