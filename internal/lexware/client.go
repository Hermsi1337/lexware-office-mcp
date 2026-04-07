package lexware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	baseURL          string
	apiToken         string
	userAgent        string
	restClient       *resty.Client
	minInterval      time.Duration
	finalizeInvoices bool

	mu          sync.Mutex
	lastRequest time.Time
}

type Request struct {
	Method string
	Path   string
	Query  map[string]string
	Body   any
	Accept string
}

type Response struct {
	StatusCode  int               `json:"statusCode"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	Body        any               `json:"body,omitempty"`
	BodyText    string            `json:"bodyText,omitempty"`
	BodyBase64  string            `json:"bodyBase64,omitempty"`
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
		baseURL:          cfg.BaseURL,
		apiToken:         cfg.APIToken,
		userAgent:        cfg.UserAgent,
		restClient:       restClient,
		minInterval:      cfg.MinInterval,
		finalizeInvoices: cfg.FinalizeInvoices,
	}
}

func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodGet
	}

	if err := c.waitTurn(ctx); err != nil {
		return nil, err
	}

	restReq := c.restClient.R().SetContext(ctx)
	if req.Body != nil {
		restReq.SetBody(req.Body)
	}
	if strings.TrimSpace(req.Accept) != "" {
		restReq.SetHeader("Accept", req.Accept)
	}
	for key, value := range req.Query {
		if strings.TrimSpace(key) == "" {
			continue
		}
		restReq.SetQueryParam(key, value)
	}

	restResp, err := restReq.Execute(method, req.Path)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	resp, err := decodeResponse(restResp.RawResponse, restResp.Body())
	if err != nil {
		return nil, err
	}

	if restResp.StatusCode() >= 400 {
		return resp, fmt.Errorf("lexware api returned status %d", restResp.StatusCode())
	}

	return resp, nil
}

func (c *Client) waitTurn(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.minInterval <= 0 {
		c.lastRequest = time.Now()
		return nil
	}

	if !c.lastRequest.IsZero() {
		wait := c.lastRequest.Add(c.minInterval).Sub(time.Now())
		if wait > 0 {
			timer := time.NewTimer(wait)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
			}
		}
	}

	c.lastRequest = time.Now()
	return nil
}

func decodeResponse(httpResp *http.Response, body []byte) (*Response, error) {
	if httpResp == nil {
		return nil, fmt.Errorf("missing http response")
	}

	resp := &Response{
		StatusCode:  httpResp.StatusCode,
		ContentType: httpResp.Header.Get("Content-Type"),
		Headers:     map[string]string{},
	}
	for key, values := range httpResp.Header {
		resp.Headers[key] = strings.Join(values, ", ")
	}

	if len(body) == 0 {
		return resp, nil
	}

	mediaType, _, _ := mime.ParseMediaType(resp.ContentType)
	switch {
	case strings.Contains(mediaType, "json") || json.Valid(body):
		var parsed any
		if err := json.Unmarshal(body, &parsed); err != nil {
			resp.BodyText = string(body)
			return resp, nil
		}
		resp.Body = parsed
	case strings.HasPrefix(mediaType, "text/"):
		resp.BodyText = string(body)
	default:
		resp.BodyBase64 = base64.StdEncoding.EncodeToString(body)
	}

	return resp, nil
}

func (c *Client) FinalizeInvoices() bool {
	return c.finalizeInvoices
}
