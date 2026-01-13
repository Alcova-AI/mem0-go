package mem0

import (
	"context"
	"net/http"
)

type SearchRequest struct {
	Query          string   `json:"query"`
	Filters        Filters  `json:"filters"`
	Version        string   `json:"version,omitempty"`
	TopK           int      `json:"top_k,omitempty"`
	Threshold      float64  `json:"threshold,omitempty"`
	Rerank         bool     `json:"rerank,omitempty"`
	KeywordSearch  bool     `json:"keyword_search,omitempty"`
	FilterMemories bool     `json:"filter_memories,omitempty"`
	Fields         []string `json:"fields,omitempty"`
	OrgID          string   `json:"org_id,omitempty"`
	ProjectID      string   `json:"project_id,omitempty"`
}

type SearchResponse struct {
	Results []Memory `json:"results"`
}

// Search performs a semantic search across memories using the given query and filters.
func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if req == nil || req.Query == "" {
		return nil, ErrMissingQuery
	}
	if req.Filters == nil {
		return nil, ErrMissingFilters
	}

	if req.Version == "" {
		req.Version = "v2"
	}
	if req.OrgID == "" && c.orgID != "" {
		req.OrgID = c.orgID
	}
	if req.ProjectID == "" && c.projectID != "" {
		req.ProjectID = c.projectID
	}

	var results []Memory
	if err := c.do(ctx, http.MethodPost, "/v2/memories/search/", nil, req, &results); err != nil {
		return nil, err
	}

	resp := &SearchResponse{Results: results}

	return resp, nil
}

// SearchUserMemories is a convenience method to search memories for a specific user.
func (c *Client) SearchUserMemories(ctx context.Context, userID, query string, opts ...SearchOption) (*SearchResponse, error) {
	req := &SearchRequest{
		Query:   query,
		Filters: NewFilters().WithUserID(userID),
	}
	for _, opt := range opts {
		opt(req)
	}
	return c.Search(ctx, req)
}

type SearchOption func(*SearchRequest)

func WithTopK(k int) SearchOption {
	return func(r *SearchRequest) { r.TopK = k }
}

func WithThreshold(t float64) SearchOption {
	return func(r *SearchRequest) { r.Threshold = t }
}

func WithRerank(enabled bool) SearchOption {
	return func(r *SearchRequest) { r.Rerank = enabled }
}

func WithKeywordSearch(enabled bool) SearchOption {
	return func(r *SearchRequest) { r.KeywordSearch = enabled }
}

func WithFilterMemories(enabled bool) SearchOption {
	return func(r *SearchRequest) { r.FilterMemories = enabled }
}

func WithFields(fields ...string) SearchOption {
	return func(r *SearchRequest) { r.Fields = fields }
}

func WithSearchFilters(filters Filters) SearchOption {
	return func(r *SearchRequest) {
		if r.Filters == nil {
			r.Filters = NewFilters()
		}
		for k, v := range filters {
			r.Filters[k] = v
		}
	}
}
