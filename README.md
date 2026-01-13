# mem0-go

A Go client for the [Mem0 Platform API](https://docs.mem0.ai/api-reference).

## Installation

```bash
go get github.com/alcova-ai/mem0-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/alcova-ai/mem0-go"
)

func main() {
    client, err := mem0.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Add a memory
    resp, err := client.AddMemory(ctx, "User prefers dark mode",
        mem0.WithUserID("user-123"),
        mem0.WithMetadata(map[string]any{"source": "settings"}),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Added memory: %s\n", resp.Results[0].ID)

    // Search memories
    searchResp, err := client.SearchUserMemories(ctx, "user-123", "preferences",
        mem0.WithTopK(5),
        mem0.WithRerank(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    for _, m := range searchResp.Results {
        fmt.Printf("Found: %s\n", m.Memory)
    }
}
```

## Features

### Memory Operations

```go
// Add memories with messages
client.AddMemories(ctx, &mem0.AddMemoriesRequest{
    Messages: []mem0.Message{
        {Role: "user", Content: "I love hiking in the mountains"},
        {Role: "assistant", Content: "That's great! Do you have a favorite trail?"},
    },
    UserID:      "user-123",
    EnableGraph: true, // Enable graph memory
})

// Get a single memory
memory, _ := client.GetMemory(ctx, "memory-id")

// Get all memories for a user
memories, _ := client.GetUserMemories(ctx, "user-123")

// Get memories with filters
memories, _ := client.GetMemories(ctx, &mem0.GetMemoriesRequest{
    Filters: mem0.NewFilters().
        WithUserID("user-123").
        WithCategories("preferences", "personal"),
    Page:     1,
    PageSize: 50,
})

// Update a memory
client.UpdateMemory(ctx, "memory-id", &mem0.UpdateMemoryRequest{
    Text: "Updated memory content",
})

// Delete a memory
client.DeleteMemory(ctx, "memory-id")

// Delete all memories for a user
client.DeleteUserMemories(ctx, "user-123")

// Get memory history
history, _ := client.GetMemoryHistory(ctx, "memory-id")
```

### Search

```go
// Basic search
results, _ := client.Search(ctx, &mem0.SearchRequest{
    Query:   "hiking preferences",
    Filters: mem0.NewFilters().WithUserID("user-123"),
    TopK:    10,
})

// Search with options
results, _ := client.SearchUserMemories(ctx, "user-123", "travel",
    mem0.WithTopK(5),
    mem0.WithThreshold(0.7),
    mem0.WithRerank(true),
    mem0.WithKeywordSearch(true),
)
```

### Advanced Filters

```go
// Complex filters with AND/OR
filters := mem0.NewFilters().
    WithUserID("user-123").
    WithCategories("travel", "food").
    WithCreatedAfter(time.Now().AddDate(0, -1, 0))

// OR filters
filters := mem0.NewFilters().WithUserID("user-123").Or(
    mem0.NewFilters().WithAgentID("agent-456"),
)
```

### Batch Operations

```go
// Batch update
client.BatchUpdate(ctx, &mem0.BatchUpdateRequest{
    Memories: []mem0.BatchUpdateItem{
        {MemoryID: "id-1", Text: "Updated text 1"},
        {MemoryID: "id-2", Text: "Updated text 2"},
    },
})

// Batch delete
client.BatchDelete(ctx, &mem0.BatchDeleteRequest{
    MemoryIDs: []string{"id-1", "id-2", "id-3"},
})
```

### Entity Management

```go
// List all users
users, _ := client.ListUsers(ctx)

// List all agents
agents, _ := client.ListAgents(ctx)

// Delete a user and all their memories
client.DeleteUser(ctx, "user-123")
```

## Client Options

```go
client, _ := mem0.NewClient("api-key",
    mem0.WithBaseURL("https://custom.api.endpoint"),
    mem0.WithTimeout(60*time.Second),
    mem0.WithOrgID("org-123"),
    mem0.WithProjectID("project-456"),
)
```

## Error Handling

```go
import "errors"

resp, err := client.GetMemory(ctx, "memory-id")
if err != nil {
    var apiErr *mem0.APIError
    if errors.As(err, &apiErr) {
        if apiErr.IsNotFound() {
            // Handle not found
        }
        if apiErr.IsRateLimited() {
            // Handle rate limit
        }
        fmt.Printf("API error: %s (status %d)\n", apiErr.Message, apiErr.StatusCode)
    }
}
```

## License

MIT
