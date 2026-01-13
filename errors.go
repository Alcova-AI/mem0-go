package mem0

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMissingAPIKey  = errors.New("mem0: API key is required")
	ErrMissingQuery   = errors.New("mem0: query is required")
	ErrMissingID      = errors.New("mem0: id is required")
	ErrMissingFilters = errors.New("mem0: filters are required")
	ErrEmptyRequest   = errors.New("mem0: request cannot be empty")
)

type APIError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type,omitempty"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
	Detail     string `json:"detail,omitempty"`
	RawBody    []byte `json:"-"`
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Detail != "" {
		return fmt.Sprintf("mem0: %s (status %d)", e.Detail, e.StatusCode)
	}
	if e.Message != "" {
		return fmt.Sprintf("mem0: %s (status %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("mem0: %s (status %d)", http.StatusText(e.StatusCode), e.StatusCode)
}

func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}
