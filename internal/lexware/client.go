package lexware

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseURL          string
	apiToken         string
	userAgent        string
	httpClient       *http.Client
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
	return &Client{
		baseURL:          cfg.BaseURL,
		apiToken:         cfg.APIToken,
		userAgent:        cfg.UserAgent,
		minInterval:      cfg.MinInterval,
		finalizeInvoices: cfg.FinalizeInvoices,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}
}

func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodGet
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = path.Join(u.Path, strings.TrimPrefix(req.Path, "/"))

	query := u.Query()
	for key, value := range req.Query {
		if strings.TrimSpace(key) == "" {
			continue
		}
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()

	var bodyReader io.Reader
	if req.Body != nil {
		payload, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiToken)
	httpReq.Header.Set("User-Agent", c.userAgent)
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	if strings.TrimSpace(req.Accept) != "" {
		httpReq.Header.Set("Accept", req.Accept)
	} else {
		httpReq.Header.Set("Accept", "application/json")
	}

	if err := c.waitTurn(ctx); err != nil {
		return nil, err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer httpResp.Body.Close()

	resp, err := decodeResponse(httpResp)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode >= 400 {
		return resp, fmt.Errorf("lexware api returned status %d", httpResp.StatusCode)
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

func decodeResponse(httpResp *http.Response) (*Response, error) {
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
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
