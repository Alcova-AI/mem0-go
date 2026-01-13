// Basic example demonstrating core memory operations
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	mem0 "github.com/alcova-ai/mem0-go"
)

func main() {
	apiKey := os.Getenv("MEM0_API_KEY")
	if apiKey == "" {
		log.Fatal("MEM0_API_KEY environment variable is required")
	}

	client, err := mem0.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	userID := "demo-user-001"

	// Add a memory
	fmt.Println("=== Adding Memory ===")
	addResp, err := client.AddMemory(ctx, "I prefer dark mode and use VS Code as my editor",
		mem0.WithUserID(userID),
		mem0.WithMetadata(map[string]any{
			"source":   "cli-demo",
			"category": "preferences",
		}),
	)
	if err != nil {
		log.Fatalf("Failed to add memory: %v", err)
	}
	fmt.Printf("Added %d memory event(s)\n", len(addResp.Results))
	for _, evt := range addResp.Results {
		id := evt.ID
		if id == "" {
			id = evt.EventID
		}
		status := evt.Event
		if status == "" {
			status = evt.Status
		}
		fmt.Printf("  - ID: %s, Status: %s\n", id, status)
	}

	// Search memories
	fmt.Println("\n=== Searching Memories ===")
	searchResp, err := client.SearchUserMemories(ctx, userID, "editor preferences",
		mem0.WithTopK(5),
	)
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}
	fmt.Printf("Found %d memories\n", len(searchResp.Results))
	for _, m := range searchResp.Results {
		fmt.Printf("  - [%s] %s (score: %.3f)\n", m.ID, m.Memory, m.Score)
	}

	// Get all memories for user
	fmt.Println("\n=== Getting All User Memories ===")
	memories, err := client.GetUserMemories(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get memories: %v", err)
	}
	fmt.Printf("User has %d memories\n", len(memories.Results))
	for _, m := range memories.Results {
		fmt.Printf("  - [%s] %s\n", m.ID, m.Memory)
	}

	fmt.Println("\n=== Done ===")
}
