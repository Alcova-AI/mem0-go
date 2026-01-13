package mem0

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type ListEntitiesRequest struct {
	Type      EntityType
	Page      int
	PageSize  int
	OrgID     string
	ProjectID string
}

type ListEntitiesResponse struct {
	Results  []Entity `json:"results"`
	Page     int      `json:"page,omitempty"`
	PageSize int      `json:"page_size,omitempty"`
	Total    int      `json:"total,omitempty"`
}

// ListEntities retrieves entities (users, agents, apps, runs) with optional filtering.
func (c *Client) ListEntities(ctx context.Context, req *ListEntitiesRequest) (*ListEntitiesResponse, error) {
	query := url.Values{}

	if req != nil {
		if req.Type != "" {
			query.Set("type", string(req.Type))
		}
		if req.Page > 0 {
			query.Set("page", strconv.Itoa(req.Page))
		}
		if req.PageSize > 0 {
			query.Set("page_size", strconv.Itoa(req.PageSize))
		}
		if req.OrgID != "" {
			query.Set("org_id", req.OrgID)
		} else if c.orgID != "" {
			query.Set("org_id", c.orgID)
		}
		if req.ProjectID != "" {
			query.Set("project_id", req.ProjectID)
		} else if c.projectID != "" {
			query.Set("project_id", c.projectID)
		}
	} else {
		if c.orgID != "" {
			query.Set("org_id", c.orgID)
		}
		if c.projectID != "" {
			query.Set("project_id", c.projectID)
		}
	}

	var resp ListEntitiesResponse
	if err := c.do(ctx, http.MethodGet, "/v1/entities/", query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListUsers retrieves all user entities.
func (c *Client) ListUsers(ctx context.Context) (*ListEntitiesResponse, error) {
	return c.ListEntities(ctx, &ListEntitiesRequest{Type: EntityTypeUser})
}

// ListAgents retrieves all agent entities.
func (c *Client) ListAgents(ctx context.Context) (*ListEntitiesResponse, error) {
	return c.ListEntities(ctx, &ListEntitiesRequest{Type: EntityTypeAgent})
}

// DeleteEntity deletes an entity and all its associated memories.
func (c *Client) DeleteEntity(ctx context.Context, entityType EntityType, entityID string) error {
	if entityType == "" || entityID == "" {
		return ErrMissingID
	}

	path := "/v2/entities/" + url.PathEscape(string(entityType)) + "/" + url.PathEscape(entityID) + "/"
	return c.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// DeleteUser deletes a user entity and all its associated memories.
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	return c.DeleteEntity(ctx, EntityTypeUser, userID)
}

// DeleteAgent deletes an agent entity and all its associated memories.
func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	return c.DeleteEntity(ctx, EntityTypeAgent, agentID)
}
