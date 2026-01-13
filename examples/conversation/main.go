// Conversation example demonstrating multi-turn memory extraction
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
	userID := "demo-user-conversation"

	// Simulate a multi-turn conversation
	conversation := []mem0.Message{
		{Role: "user", Content: "Hi! I'm planning a trip to Japan next month."},
		{Role: "assistant", Content: "That sounds exciting! Is this your first time visiting Japan?"},
		{Role: "user", Content: "Yes, it's my first time! I'm really interested in visiting temples in Kyoto and trying authentic ramen."},
		{Role: "assistant", Content: "Great choices! Kyoto has beautiful temples like Kinkaku-ji and Fushimi Inari. For ramen, I'd recommend trying Ichiran or local spots in Osaka."},
		{Role: "user", Content: "Perfect! I'm vegetarian though, so I'll need to find places with veggie options."},
	}

	fmt.Println("=== Adding Conversation ===")
	fmt.Println("Conversation:")
	for _, msg := range conversation {
		fmt.Printf("  [%s]: %s\n", msg.Role, msg.Content)
	}

	// Add the conversation - Mem0 will extract relevant facts
	resp, err := client.AddMemories(ctx, &mem0.AddMemoriesRequest{
		Messages: conversation,
		UserID:   userID,
		Metadata: map[string]any{
			"conversation_id": "conv-001",
			"source":          "cli-demo",
		},
	})
	if err != nil {
		log.Fatalf("Failed to add conversation: %v", err)
	}

	fmt.Printf("\nExtracted %d memory events:\n", len(resp.Results))
	for _, evt := range resp.Results {
		status := evt.Event
		if status == "" {
			status = evt.Status
		}
		memory := evt.Memory
		if memory == "" {
			memory = evt.Message
		}
		fmt.Printf("  - %s: %s\n", status, memory)
	}

	// Now search for what we learned
	fmt.Println("\n=== Searching for Travel Plans ===")
	searchResp, err := client.SearchUserMemories(ctx, userID, "travel dietary preferences")
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}

	fmt.Printf("Found %d relevant memories:\n", len(searchResp.Results))
	for _, m := range searchResp.Results {
		fmt.Printf("  - %s\n", m.Memory)
	}
}
