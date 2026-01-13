# mem0-go Examples

CLI examples demonstrating the mem0-go client library.

## Prerequisites

1. Get an API key from [Mem0 Dashboard](https://app.mem0.ai/dashboard/api-keys)
2. Set the environment variable:
   ```bash
   export MEM0_API_KEY="your-api-key"
   ```

## Running Examples

From the `mem0-go` directory:

### Basic Usage
Demonstrates adding, searching, and retrieving memories:
```bash
go run ./examples/basic
```

### Conversation
Shows how to extract memories from multi-turn conversations:
```bash
go run ./examples/conversation
```

### Search
Advanced search features including reranking, thresholds, and keyword search:
```bash
go run ./examples/search
```

### CRUD Operations
Full create, read, update, delete lifecycle with history:
```bash
go run ./examples/crud
```

### Batch Operations
Batch update and delete operations:
```bash
go run ./examples/batch
```

### Entities
User and agent management:
```bash
go run ./examples/entities
```

### Wealth Advisor
Real-world example of a wealth advisor conversation with life event extraction:
```bash
go run ./examples/wealth-advisor
```

This example demonstrates:
- Multi-turn advisor/client conversation
- Automatic extraction of life events (house purchase, family planning)
- Financial details extraction (mortgage, 401k, savings)
- Semantic search for relevant context in future sessions

## Example Output

### Basic Example
```
=== Adding Memory ===
Added 1 memory event(s)
  - ID: abc123, Event: ADD

=== Searching Memories ===
Found 1 memories
  - [abc123] User prefers dark mode (score: 0.892)

=== Getting All User Memories ===
User has 1 memories
  - [abc123] User prefers dark mode
```

### Conversation Example
```
=== Adding Conversation ===
Conversation:
  [user]: Hi! I'm planning a trip to Japan next month.
  [assistant]: That sounds exciting! Is this your first time visiting Japan?
  ...

Extracted 3 memory events:
  - ADD: User is planning a trip to Japan
  - ADD: User is interested in temples in Kyoto
  - ADD: User is vegetarian
```

## Tips

- All examples use unique user IDs to avoid conflicts
- Examples clean up after themselves where possible
- Use `MEM0_API_KEY` environment variable for authentication
- For production, consider using `WithOrgID` and `WithProjectID` client options
