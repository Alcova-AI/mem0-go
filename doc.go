// Package mem0 provides a Go client for the mem0.ai API.
//
// The client supports memory CRUD operations, semantic search, entity management,
// and batch operations for efficient memory management.
//
// # Quick Start
//
//	client, err := mem0.NewClient("your-api-key")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Add a memory
//	resp, err := client.AddMemory(ctx, "User prefers dark mode",
//	    mem0.WithUserID("user-123"),
//	)
//
//	// Search memories
//	results, err := client.SearchUserMemories(ctx, "user-123", "preferences")
//
// # Filters
//
// Use [Filters] to build complex queries:
//
//	filters := mem0.NewFilters().
//	    WithUserID("user-123").
//	    WithCategories("preferences", "settings").
//	    WithCreatedAfter(time.Now().AddDate(0, -1, 0))
//
// Filters can be combined with And/Or:
//
//	combined := mem0.NewFilters().WithUserID("user-1").
//	    Or(mem0.NewFilters().WithUserID("user-2"))
//
// # Configuration
//
// The client supports various options:
//
//	client, err := mem0.NewClient("api-key",
//	    mem0.WithBaseURL("https://custom.api"),
//	    mem0.WithTimeout(60*time.Second),
//	    mem0.WithOrgID("org-123"),
//	    mem0.WithProjectID("proj-456"),
//	)
package mem0
