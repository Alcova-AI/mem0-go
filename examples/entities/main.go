// Entities example demonstrating user/agent management
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

	// Create some memories for different entity types
	fmt.Println("=== Setting Up Test Entities ===")

	// User memory
	_, err = client.AddMemory(ctx, "User prefers email communication",
		mem0.WithUserID("entity-demo-user"),
	)
	if err != nil {
		log.Printf("Warning: %v", err)
	}
	fmt.Println("Added memory for user: entity-demo-user")

	// Agent memory
	_, err = client.AddMemories(ctx, &mem0.AddMemoriesRequest{
		Messages: []mem0.Message{
			{Role: "system", Content: "This agent specializes in customer support"},
		},
		AgentID: "entity-demo-agent",
	})
	if err != nil {
		log.Printf("Warning: %v", err)
	}
	fmt.Println("Added memory for agent: entity-demo-agent")

	// List all users
	fmt.Println("\n=== List Users ===")
	users, err := client.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	fmt.Printf("Found %d users:\n", len(users.Results))
	for _, u := range users.Results {
		fmt.Printf("  - %s (type: %s, memories: %d)\n", u.ID, u.Type, u.TotalMemories)
	}

	// List all agents
	fmt.Println("\n=== List Agents ===")
	agents, err := client.ListAgents(ctx)
	if err != nil {
		log.Fatalf("Failed to list agents: %v", err)
	}
	fmt.Printf("Found %d agents:\n", len(agents.Results))
	for _, a := range agents.Results {
		fmt.Printf("  - %s (type: %s, memories: %d)\n", a.ID, a.Type, a.TotalMemories)
	}

	// List all entities with pagination
	fmt.Println("\n=== List All Entities (Paginated) ===")
	entities, err := client.ListEntities(ctx, &mem0.ListEntitiesRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("Failed to list entities: %v", err)
	}
	fmt.Printf("Found %d entities (page 1):\n", len(entities.Results))
	for _, e := range entities.Results {
		fmt.Printf("  - %s [%s]: %d memories\n", e.ID, e.Type, e.TotalMemories)
	}

	// Cleanup - delete the test entities
	fmt.Println("\n=== Cleanup ===")

	err = client.DeleteUser(ctx, "entity-demo-user")
	if err != nil {
		log.Printf("Warning deleting user: %v", err)
	} else {
		fmt.Println("Deleted user: entity-demo-user")
	}

	err = client.DeleteAgent(ctx, "entity-demo-agent")
	if err != nil {
		log.Printf("Warning deleting agent: %v", err)
	} else {
		fmt.Println("Deleted agent: entity-demo-agent")
	}
}
