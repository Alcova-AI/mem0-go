package mem0

import (
	"context"
	"net/http"
)

type AddMemoriesRequest struct {
	Messages           []Message      `json:"messages,omitempty"`
	UserID             string         `json:"user_id,omitempty"`
	AgentID            string         `json:"agent_id,omitempty"`
	AppID              string         `json:"app_id,omitempty"`
	RunID              string         `json:"run_id,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
	Infer              *bool          `json:"infer,omitempty"`
	Immutable          bool           `json:"immutable,omitempty"`
	ExpirationDate     string         `json:"expiration_date,omitempty"`
	CustomCategories   map[string]any `json:"custom_categories,omitempty"`
	CustomInstructions string         `json:"custom_instructions,omitempty"`
	Includes           string         `json:"includes,omitempty"`
	Excludes           string         `json:"excludes,omitempty"`
	EnableGraph        bool           `json:"enable_graph,omitempty"`
	AsyncMode          *bool          `json:"async_mode,omitempty"`
	OutputFormat       string         `json:"output_format,omitempty"`
	OrgID              string         `json:"org_id,omitempty"`
	ProjectID          string         `json:"project_id,omitempty"`
	Timestamp          int64          `json:"timestamp,omitempty"` // Unix timestamp for temporal context
}

type AddMemoriesResponse struct {
	Results []AddEvent `json:"results"`
}

func (c *Client) AddMemories(ctx context.Context, req *AddMemoriesRequest) (*AddMemoriesResponse, error) {
	if req == nil || len(req.Messages) == 0 {
		return nil, ErrEmptyRequest
	}

	if req.OutputFormat == "" {
		req.OutputFormat = "v1.1"
	}
	if req.OrgID == "" && c.orgID != "" {
		req.OrgID = c.orgID
	}
	if req.ProjectID == "" && c.projectID != "" {
		req.ProjectID = c.projectID
	}

	var resp AddMemoriesResponse
	if err := c.do(ctx, http.MethodPost, "/v1/memories/", nil, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) AddMemory(ctx context.Context, content string, opts ...AddMemoryOption) (*AddMemoriesResponse, error) {
	req := &AddMemoriesRequest{
		Messages: []Message{{Role: "user", Content: content}},
	}
	for _, opt := range opts {
		opt(req)
	}
	return c.AddMemories(ctx, req)
}

type AddMemoryOption func(*AddMemoriesRequest)

func WithUserID(id string) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.UserID = id }
}

func WithAgentID(id string) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.AgentID = id }
}

func WithAppID(id string) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.AppID = id }
}

func WithRunID(id string) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.RunID = id }
}

func WithMetadata(m map[string]any) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.Metadata = m }
}

func WithGraph(enabled bool) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.EnableGraph = enabled }
}

func WithImmutable(immutable bool) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.Immutable = immutable }
}

func WithExpiration(date string) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.ExpirationDate = date }
}

func WithInfer(infer bool) AddMemoryOption {
	return func(r *AddMemoriesRequest) { r.Infer = &infer }
}

func (c *Client) GetMemory(ctx context.Context, memoryID string) (*Memory, error) {
	if memoryID == "" {
		return nil, ErrMissingID
	}

	var mem Memory
	if err := c.do(ctx, http.MethodGet, "/v1/memories/"+memoryID+"/", nil, nil, &mem); err != nil {
		return nil, err
	}

	return &mem, nil
}

type GetMemoriesRequest struct {
	Filters   Filters  `json:"filters"`
	Fields    []string `json:"fields,omitempty"`
	Page      int      `json:"page,omitempty"`
	PageSize  int      `json:"page_size,omitempty"`
	OrgID     string   `json:"org_id,omitempty"`
	ProjectID string   `json:"project_id,omitempty"`
}

type GetMemoriesResponse struct {
	Results  []Memory `json:"results"`
	Page     int      `json:"page,omitempty"`
	PageSize int      `json:"page_size,omitempty"`
	Total    int      `json:"total,omitempty"`
}

func (c *Client) GetMemories(ctx context.Context, req *GetMemoriesRequest) (*GetMemoriesResponse, error) {
	if req == nil || req.Filters == nil {
		return nil, ErrMissingFilters
	}

	if req.OrgID == "" && c.orgID != "" {
		req.OrgID = c.orgID
	}
	if req.ProjectID == "" && c.projectID != "" {
		req.ProjectID = c.projectID
	}

	var results []Memory
	if err := c.do(ctx, http.MethodPost, "/v2/memories/", nil, req, &results); err != nil {
		return nil, err
	}

	return &GetMemoriesResponse{Results: results}, nil
}

func (c *Client) GetUserMemories(ctx context.Context, userID string) (*GetMemoriesResponse, error) {
	return c.GetMemories(ctx, &GetMemoriesRequest{
		Filters: NewFilters().WithUserID(userID),
	})
}

type UpdateMemoryRequest struct {
	Text     string         `json:"text,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// UpdateMemory updates an existing memory's text or metadata.
func (c *Client) UpdateMemory(ctx context.Context, memoryID string, req *UpdateMemoryRequest) (*Memory, error) {
	if memoryID == "" {
		return nil, ErrMissingID
	}
	if req == nil {
		return nil, ErrEmptyRequest
	}

	var mem Memory
	if err := c.do(ctx, http.MethodPut, "/v1/memories/"+memoryID+"/", nil, req, &mem); err != nil {
		return nil, err
	}

	return &mem, nil
}

func (c *Client) DeleteMemory(ctx context.Context, memoryID string) error {
	if memoryID == "" {
		return ErrMissingID
	}

	return c.do(ctx, http.MethodDelete, "/v1/memories/"+memoryID+"/", nil, nil, nil)
}

type DeleteMemoriesRequest struct {
	Filters   Filters `json:"filters"`
	OrgID     string  `json:"org_id,omitempty"`
	ProjectID string  `json:"project_id,omitempty"`
}

func (c *Client) DeleteMemories(ctx context.Context, req *DeleteMemoriesRequest) error {
	if req == nil || req.Filters == nil {
		return ErrMissingFilters
	}

	if req.OrgID == "" && c.orgID != "" {
		req.OrgID = c.orgID
	}
	if req.ProjectID == "" && c.projectID != "" {
		req.ProjectID = c.projectID
	}

	return c.do(ctx, http.MethodDelete, "/v1/memories/all/", nil, req, nil)
}

func (c *Client) DeleteUserMemories(ctx context.Context, userID string) error {
	return c.DeleteMemories(ctx, &DeleteMemoriesRequest{
		Filters: NewFilters().WithUserID(userID),
	})
}

// GetMemoryHistory retrieves the change history for a memory.
func (c *Client) GetMemoryHistory(ctx context.Context, memoryID string) ([]MemoryHistory, error) {
	if memoryID == "" {
		return nil, ErrMissingID
	}

	var history []MemoryHistory
	if err := c.do(ctx, http.MethodGet, "/v1/memories/"+memoryID+"/history/", nil, nil, &history); err != nil {
		return nil, err
	}

	return history, nil
}

type BatchUpdateItem struct {
	MemoryID string         `json:"memory_id"`
	Text     string         `json:"text,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type BatchUpdateRequest struct {
	Memories []BatchUpdateItem `json:"memories"`
}

type BatchUpdateResponse struct {
	Message string `json:"message"`
}

func (c *Client) BatchUpdate(ctx context.Context, req *BatchUpdateRequest) (*BatchUpdateResponse, error) {
	if req == nil || len(req.Memories) == 0 {
		return nil, ErrEmptyRequest
	}

	var resp BatchUpdateResponse
	if err := c.do(ctx, http.MethodPut, "/v1/batch/", nil, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type BatchDeleteRequest struct {
	MemoryIDs []string
}

type batchDeleteItem struct {
	MemoryID string `json:"memory_id"`
}

type batchDeleteBody struct {
	Memories []batchDeleteItem `json:"memories"`
}

// BatchDelete deletes multiple memories in a single request.
func (c *Client) BatchDelete(ctx context.Context, req *BatchDeleteRequest) error {
	if req == nil || len(req.MemoryIDs) == 0 {
		return ErrEmptyRequest
	}

	body := batchDeleteBody{
		Memories: make([]batchDeleteItem, len(req.MemoryIDs)),
	}
	for i, id := range req.MemoryIDs {
		body.Memories[i] = batchDeleteItem{MemoryID: id}
	}

	return c.do(ctx, http.MethodDelete, "/v1/batch/", nil, body, nil)
}
