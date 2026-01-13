package mem0

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL   = "https://api.mem0.ai"
	defaultUserAgent = "mem0-go/0.1.0"
	defaultTimeout   = 30 * time.Second
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	userAgent  string
	orgID      string
	projectID  string
}

func NewClient(apiKey string, opts ...ClientOption) (*Client, error) {
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	c := &Client{
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		userAgent: defaultUserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("mem0: invalid URL: %w", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("mem0: failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("mem0: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mem0: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mem0: failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    respBody,
		}
		_ = json.Unmarshal(respBody, apiErr)
		return apiErr
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("mem0: failed to unmarshal response: %w", err)
		}
	}

	return nil
}
