// CRUD example demonstrating full lifecycle of memories
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	userID := "demo-user-crud"

	// CREATE (async mode - queued for background processing)
	fmt.Println("=== CREATE ===")
	addResp, err := client.AddMemory(ctx, "My favorite color is blue",
		mem0.WithUserID(userID),
	)
	if err != nil {
		log.Fatalf("Failed to create: %v", err)
	}
	if len(addResp.Results) == 0 {
		log.Fatal("No memory created")
	}
	eventID := addResp.Results[0].EventID
	fmt.Printf("Memory queued with event ID: %s (Status: %s)\n", eventID, addResp.Results[0].Status)

	// Wait for async processing with retries
	fmt.Println("Waiting for memory to be processed...")
	var memoryID string
	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)
		searchResp, err := client.SearchUserMemories(ctx, userID, "favorite color")
		if err != nil {
			log.Fatalf("Failed to search: %v", err)
		}
		if len(searchResp.Results) > 0 {
			memoryID = searchResp.Results[0].ID
			break
		}
		fmt.Printf("  Retry %d/5...\n", i+1)
	}
	if memoryID == "" {
		log.Fatal("Memory not found after waiting - may need more time")
	}
	fmt.Printf("Found memory with ID: %s\n", memoryID)

	// READ
	fmt.Println("\n=== READ ===")
	memory, err := client.GetMemory(ctx, memoryID)
	if err != nil {
		log.Fatalf("Failed to read: %v", err)
	}
	fmt.Printf("Memory: %s\n", memory.Memory)
	fmt.Printf("User ID: %s\n", memory.UserID)
	fmt.Printf("Created: %s\n", memory.CreatedAt)

	// UPDATE
	fmt.Println("\n=== UPDATE ===")
	updated, err := client.UpdateMemory(ctx, memoryID, &mem0.UpdateMemoryRequest{
		Text: "My favorite color is green (changed from blue)",
	})
	if err != nil {
		log.Fatalf("Failed to update: %v", err)
	}
	fmt.Printf("Updated memory: %s\n", updated.Memory)

	// HISTORY
	fmt.Println("\n=== HISTORY ===")
	history, err := client.GetMemoryHistory(ctx, memoryID)
	if err != nil {
		log.Fatalf("Failed to get history: %v", err)
	}
	fmt.Printf("Memory has %d history entries:\n", len(history))
	for _, h := range history {
		fmt.Printf("  - Event: %s\n", h.Event)
		if h.OldMemory != "" {
			fmt.Printf("    Old: %s\n", h.OldMemory)
		}
		fmt.Printf("    New: %s\n", h.NewMemory)
	}

	// DELETE
	fmt.Println("\n=== DELETE ===")
	err = client.DeleteMemory(ctx, memoryID)
	if err != nil {
		log.Fatalf("Failed to delete: %v", err)
	}
	fmt.Println("Memory deleted successfully")

	// Verify deletion
	_, err = client.GetMemory(ctx, memoryID)
	if err != nil {
		if apiErr, ok := err.(*mem0.APIError); ok && apiErr.IsNotFound() {
			fmt.Println("Verified: memory no longer exists")
		} else {
			fmt.Printf("Unexpected error: %v\n", err)
		}
	}
}
