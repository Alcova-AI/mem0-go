package mem0

import (
	"net/http"
	"time"
)

type ClientOption func(*Client)

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

func WithOrgID(orgID string) ClientOption {
	return func(c *Client) {
		c.orgID = orgID
	}
}

func WithProjectID(projectID string) ClientOption {
	return func(c *Client) {
		c.projectID = projectID
	}
}
