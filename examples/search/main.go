// Search example demonstrating advanced search features
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
	userID := "demo-user-search"

	// First, add some test memories
	fmt.Println("=== Setting Up Test Data ===")
	testMemories := []string{
		"I work as a software engineer at a tech startup",
		"My favorite programming language is Go",
		"I enjoy hiking on weekends",
		"I'm learning Japanese in my free time",
		"I prefer tea over coffee",
	}

	for _, content := range testMemories {
		_, err := client.AddMemory(ctx, content, mem0.WithUserID(userID))
		if err != nil {
			log.Printf("Warning: failed to add memory: %v", err)
		}
	}
	fmt.Printf("Added %d test memories\n", len(testMemories))

	// Basic search
	fmt.Println("\n=== Basic Search ===")
	resp, err := client.SearchUserMemories(ctx, userID, "programming")
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	printResults("Programming search", resp.Results)

	// Search with reranking
	fmt.Println("\n=== Search with Reranking ===")
	resp, err = client.SearchUserMemories(ctx, userID, "hobbies and interests",
		mem0.WithRerank(true),
		mem0.WithTopK(3),
	)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	printResults("Hobbies (reranked)", resp.Results)

	// Search with threshold
	fmt.Println("\n=== Search with High Threshold ===")
	resp, err = client.SearchUserMemories(ctx, userID, "work",
		mem0.WithThreshold(0.7),
		mem0.WithTopK(10),
	)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	printResults("Work (threshold 0.7)", resp.Results)

	// Keyword search
	fmt.Println("\n=== Keyword Search ===")
	resp, err = client.SearchUserMemories(ctx, userID, "Go",
		mem0.WithKeywordSearch(true),
	)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	printResults("Keyword 'Go'", resp.Results)

	// Advanced search with custom filters
	fmt.Println("\n=== Advanced Search with Filters ===")
	resp, err = client.Search(ctx, &mem0.SearchRequest{
		Query:   "preferences",
		Filters: mem0.NewFilters().WithUserID(userID),
		TopK:    5,
		Rerank:  true,
	})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	printResults("Preferences (advanced)", resp.Results)
}

func printResults(label string, results []mem0.Memory) {
	fmt.Printf("%s - Found %d results:\n", label, len(results))
	for i, m := range results {
		score := ""
		if m.Score > 0 {
			score = fmt.Sprintf(" (score: %.3f)", m.Score)
		}
		fmt.Printf("  %d. %s%s\n", i+1, m.Memory, score)
	}
}
