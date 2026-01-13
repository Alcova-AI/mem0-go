package mem0

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient("")
	if err != ErrMissingAPIKey {
		t.Errorf("expected ErrMissingAPIKey, got %v", err)
	}

	client, err := NewClient("test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.apiKey != "test-key" {
		t.Errorf("expected apiKey 'test-key', got %q", client.apiKey)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", defaultBaseURL, client.baseURL)
	}
}

func TestClientWithOptions(t *testing.T) {
	client, _ := NewClient("test-key",
		WithBaseURL("https://custom.api"),
		WithOrgID("org-123"),
		WithProjectID("proj-456"),
	)

	if client.baseURL != "https://custom.api" {
		t.Errorf("expected baseURL 'https://custom.api', got %q", client.baseURL)
	}
	if client.orgID != "org-123" {
		t.Errorf("expected orgID 'org-123', got %q", client.orgID)
	}
	if client.projectID != "proj-456" {
		t.Errorf("expected projectID 'proj-456', got %q", client.projectID)
	}
}

func TestAddMemory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/" {
			t.Errorf("expected /v1/memories/, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("missing Content-Type header")
		}
		if r.Header.Get("Authorization") != "Token test-key" {
			t.Errorf("missing or invalid Authorization header")
		}

		var req AddMemoriesRequest
		json.NewDecoder(r.Body).Decode(&req)

		if len(req.Messages) == 0 {
			t.Error("expected messages in request")
		}
		if req.UserID != "user-123" {
			t.Errorf("expected UserID 'user-123', got %q", req.UserID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AddMemoriesResponse{
			Results: []AddEvent{{ID: "mem-1", Event: "ADD"}},
		})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	resp, err := client.AddMemory(context.Background(), "test content", WithUserID("user-123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].ID != "mem-1" {
		t.Errorf("expected ID 'mem-1', got %q", resp.Results[0].ID)
	}
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v2/memories/search/" {
			t.Errorf("expected /v2/memories/search/, got %s", r.URL.Path)
		}

		var req SearchRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Query != "test query" {
			t.Errorf("expected query 'test query', got %q", req.Query)
		}
		if req.TopK != 5 {
			t.Errorf("expected TopK 5, got %d", req.TopK)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Memory{{ID: "mem-1", Memory: "test memory"}})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	resp, err := client.SearchUserMemories(context.Background(), "user-123", "test query", WithTopK(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}
}

func TestGetMemory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/memories/mem-123/" {
			t.Errorf("expected /v1/memories/mem-123/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Memory{
			ID:     "mem-123",
			Memory: "test memory content",
			UserID: "user-456",
		})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	mem, err := client.GetMemory(context.Background(), "mem-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mem.ID != "mem-123" {
		t.Errorf("expected ID 'mem-123', got %q", mem.ID)
	}
	if mem.Memory != "test memory content" {
		t.Errorf("expected Memory 'test memory content', got %q", mem.Memory)
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIError{
			Message: "Memory not found",
		})
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	_, err := client.GetMemory(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Errorf("expected IsNotFound() to return true")
	}
}

func TestFilters(t *testing.T) {
	f := NewFilters().
		WithUserID("user-123").
		WithCategories("travel", "food")

	if f["user_id"] != "user-123" {
		t.Errorf("expected user_id 'user-123', got %v", f["user_id"])
	}

	cats, ok := f["categories"].(map[string]any)
	if !ok {
		t.Fatalf("expected categories to be map, got %T", f["categories"])
	}
	if cats["in"] == nil {
		t.Error("expected categories.in to be set")
	}
}

func TestDeleteEntity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v2/entities/user/user-123/" {
			t.Errorf("expected /v2/entities/user/user-123/, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient("test-key", WithBaseURL(server.URL))

	err := client.DeleteUser(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.Search(context.Background(), nil)
	if err != ErrMissingQuery {
		t.Errorf("expected ErrMissingQuery for nil request, got %v", err)
	}

	_, err = client.Search(context.Background(), &SearchRequest{})
	if err != ErrMissingQuery {
		t.Errorf("expected ErrMissingQuery for empty query, got %v", err)
	}

	_, err = client.Search(context.Background(), &SearchRequest{Query: "test"})
	if err != ErrMissingFilters {
		t.Errorf("expected ErrMissingFilters for nil filters, got %v", err)
	}
}

func TestMemoryValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.GetMemory(context.Background(), "")
	if err != ErrMissingID {
		t.Errorf("GetMemory: expected ErrMissingID, got %v", err)
	}

	err = client.DeleteMemory(context.Background(), "")
	if err != ErrMissingID {
		t.Errorf("DeleteMemory: expected ErrMissingID, got %v", err)
	}

	_, err = client.UpdateMemory(context.Background(), "", nil)
	if err != ErrMissingID {
		t.Errorf("UpdateMemory: expected ErrMissingID for empty ID, got %v", err)
	}

	_, err = client.UpdateMemory(context.Background(), "mem-1", nil)
	if err != ErrEmptyRequest {
		t.Errorf("UpdateMemory: expected ErrEmptyRequest for nil request, got %v", err)
	}

	_, err = client.AddMemories(context.Background(), nil)
	if err != ErrEmptyRequest {
		t.Errorf("AddMemories: expected ErrEmptyRequest for nil request, got %v", err)
	}

	_, err = client.AddMemories(context.Background(), &AddMemoriesRequest{})
	if err != ErrEmptyRequest {
		t.Errorf("AddMemories: expected ErrEmptyRequest for empty messages, got %v", err)
	}
}

func TestBatchValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.BatchUpdate(context.Background(), nil)
	if err != ErrEmptyRequest {
		t.Errorf("BatchUpdate: expected ErrEmptyRequest for nil request, got %v", err)
	}

	_, err = client.BatchUpdate(context.Background(), &BatchUpdateRequest{})
	if err != ErrEmptyRequest {
		t.Errorf("BatchUpdate: expected ErrEmptyRequest for empty memories, got %v", err)
	}

	err = client.BatchDelete(context.Background(), nil)
	if err != ErrEmptyRequest {
		t.Errorf("BatchDelete: expected ErrEmptyRequest for nil request, got %v", err)
	}

	err = client.BatchDelete(context.Background(), &BatchDeleteRequest{})
	if err != ErrEmptyRequest {
		t.Errorf("BatchDelete: expected ErrEmptyRequest for empty IDs, got %v", err)
	}
}

func TestDeleteMemoriesValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	err := client.DeleteMemories(context.Background(), nil)
	if err != ErrMissingFilters {
		t.Errorf("expected ErrMissingFilters for nil request, got %v", err)
	}

	err = client.DeleteMemories(context.Background(), &DeleteMemoriesRequest{})
	if err != ErrMissingFilters {
		t.Errorf("expected ErrMissingFilters for nil filters, got %v", err)
	}
}

func TestGetMemoriesValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.GetMemories(context.Background(), nil)
	if err != ErrMissingFilters {
		t.Errorf("expected ErrMissingFilters for nil request, got %v", err)
	}

	_, err = client.GetMemories(context.Background(), &GetMemoriesRequest{})
	if err != ErrMissingFilters {
		t.Errorf("expected ErrMissingFilters for nil filters, got %v", err)
	}
}

func TestFiltersDateRange(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	f := NewFilters().
		WithCreatedAfter(yesterday).
		WithCreatedBefore(now)

	created, ok := f["created_at"].(map[string]any)
	if !ok {
		t.Fatalf("expected created_at to be map, got %T", f["created_at"])
	}

	if created["gte"] == nil {
		t.Error("expected created_at.gte to be set")
	}
	if created["lte"] == nil {
		t.Error("expected created_at.lte to be set")
	}
}

func TestDeleteEntityValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	err := client.DeleteEntity(context.Background(), "", "user-123")
	if err != ErrMissingID {
		t.Errorf("expected ErrMissingID for empty entity type, got %v", err)
	}

	err = client.DeleteEntity(context.Background(), EntityTypeUser, "")
	if err != ErrMissingID {
		t.Errorf("expected ErrMissingID for empty entity ID, got %v", err)
	}
}

func TestGetMemoryHistoryValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.GetMemoryHistory(context.Background(), "")
	if err != ErrMissingID {
		t.Errorf("expected ErrMissingID, got %v", err)
	}
}
