package mem0

import "time"

type Memory struct {
	ID             string         `json:"id"`
	Memory         string         `json:"memory"`
	UserID         string         `json:"user_id,omitempty"`
	AgentID        string         `json:"agent_id,omitempty"`
	AppID          string         `json:"app_id,omitempty"`
	RunID          string         `json:"run_id,omitempty"`
	Hash           string         `json:"hash,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	Categories     []string       `json:"categories,omitempty"`
	Immutable      bool           `json:"immutable,omitempty"`
	ExpirationDate string         `json:"expiration_date,omitempty"`
	Owner          string         `json:"owner,omitempty"`
	Organization   string         `json:"organization,omitempty"`
	CreatedAt      time.Time      `json:"created_at,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at,omitempty"`
	Score          float64        `json:"score,omitempty"`
}

type MemoryHistory struct {
	ID        string         `json:"id"`
	MemoryID  string         `json:"memory_id"`
	OldMemory string         `json:"old_memory,omitempty"`
	NewMemory string         `json:"new_memory"`
	Event     string         `json:"event"`
	UserID    string         `json:"user_id,omitempty"`
	Input     []Message      `json:"input,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Entity struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"`
	TotalMemories int            `json:"total_memories"`
	Owner         string         `json:"owner,omitempty"`
	Organization  string         `json:"organization,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type EntityType string

const (
	EntityTypeUser  EntityType = "user"
	EntityTypeAgent EntityType = "agent"
	EntityTypeApp   EntityType = "app"
	EntityTypeRun   EntityType = "run"
)

type AddEvent struct {
	ID      string `json:"id,omitempty"`
	EventID string `json:"event_id,omitempty"`
	Event   string `json:"event,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Memory  string `json:"memory,omitempty"`
}

type Filters map[string]any

func NewFilters() Filters {
	return make(Filters)
}

func (f Filters) WithUserID(id string) Filters {
	f["user_id"] = id
	return f
}

func (f Filters) WithAgentID(id string) Filters {
	f["agent_id"] = id
	return f
}

func (f Filters) WithAppID(id string) Filters {
	f["app_id"] = id
	return f
}

func (f Filters) WithRunID(id string) Filters {
	f["run_id"] = id
	return f
}

func (f Filters) WithCategories(categories ...string) Filters {
	f["categories"] = map[string]any{"in": categories}
	return f
}

func (f Filters) WithCategoryContains(category string) Filters {
	f["categories"] = map[string]any{"contains": category}
	return f
}

func (f Filters) WithCreatedAfter(t time.Time) Filters {
	m, _ := f["created_at"].(map[string]any)
	if m == nil {
		m = map[string]any{}
		f["created_at"] = m
	}
	m["gte"] = t.Format(time.RFC3339)
	return f
}

func (f Filters) WithCreatedBefore(t time.Time) Filters {
	m, _ := f["created_at"].(map[string]any)
	if m == nil {
		m = map[string]any{}
		f["created_at"] = m
	}
	m["lte"] = t.Format(time.RFC3339)
	return f
}

func (f Filters) And(filters ...Filters) Filters {
	parts := make([]Filters, 0, len(filters)+1)
	parts = append(parts, f)
	parts = append(parts, filters...)
	return Filters{"AND": parts}
}

func (f Filters) Or(filters ...Filters) Filters {
	parts := make([]Filters, 0, len(filters)+1)
	parts = append(parts, f)
	parts = append(parts, filters...)
	return Filters{"OR": parts}
}
