// Batch operations example
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
	userID := "demo-user-batch"

	// Create multiple memories
	fmt.Println("=== Creating Multiple Memories ===")
	memories := []string{
		"Batch memory 1: I like pizza",
		"Batch memory 2: I work from home",
		"Batch memory 3: My timezone is PST",
	}

	for _, content := range memories {
		resp, err := client.AddMemory(ctx, content, mem0.WithUserID(userID))
		if err != nil {
			log.Printf("Warning: failed to add memory: %v", err)
			continue
		}
		if len(resp.Results) > 0 {
			fmt.Printf("Queued: %s (status: %s)\n", resp.Results[0].EventID, resp.Results[0].Status)
		}
	}

	// Wait for async processing
	fmt.Println("\nWaiting for memories to be processed...")
	time.Sleep(3 * time.Second)

	// Get actual memory IDs
	allMemories, err := client.GetUserMemories(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to get memories: %v", err)
	}

	fmt.Printf("\n=== Found %d Memories ===\n", len(allMemories.Results))
	var memoryIDs []string
	for _, m := range allMemories.Results {
		fmt.Printf("  - [%s] %s\n", m.ID, m.Memory)
		memoryIDs = append(memoryIDs, m.ID)
	}

	if len(memoryIDs) < 2 {
		fmt.Println("Need at least 2 memories for batch operations, skipping batch update")
		return
	}

	// Batch update
	fmt.Println("\n=== Batch Update ===")
	updateResp, err := client.BatchUpdate(ctx, &mem0.BatchUpdateRequest{
		Memories: []mem0.BatchUpdateItem{
			{MemoryID: memoryIDs[0], Text: "Updated: I LOVE pizza"},
			{MemoryID: memoryIDs[1], Text: "Updated: I work from home 3 days a week"},
		},
	})
	if err != nil {
		log.Fatalf("Batch update failed: %v", err)
	}
	fmt.Printf("Result: %s\n", updateResp.Message)

	// Batch delete
	fmt.Println("\n=== Batch Delete ===")
	err = client.BatchDelete(ctx, &mem0.BatchDeleteRequest{
		MemoryIDs: memoryIDs,
	})
	if err != nil {
		log.Fatalf("Batch delete failed: %v", err)
	}
	fmt.Printf("Deleted %d memories\n", len(memoryIDs))

	// Verify
	fmt.Println("\n=== Verifying Deletion ===")
	remaining, err := client.GetUserMemories(ctx, userID)
	if err != nil {
		log.Fatalf("Failed to verify: %v", err)
	}
	fmt.Printf("Remaining memories for user: %d\n", len(remaining.Results))
}
