package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	var (
		serverURL = flag.String("address", "https://app.recallkitchen.com/mcp", "MCP server HTTP endpoint")
		timeout   = flag.Duration("timeout", 15*time.Second, "request timeout")
	)
	flag.Parse()

	// 1. Create client identity
	clientImpl := &mcp.Implementation{
		Name:    "rk-mcp-go-example",
		Version: "v0.1.0",
	}

	// 2. Create the client (options are optional)
	client := mcp.NewClient(clientImpl, nil)

	// 3. Connect (handshake + capability discovery)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: *serverURL}, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	fmt.Printf("Connected to MCP server (session ID: %s)\n", session.ID())

	// 4. Example: list available tools (very common first action)
	toolsCtx, toolsCancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer toolsCancel()

	toolsResult, err := session.ListTools(toolsCtx, nil)
	if err != nil {
		log.Printf("ListTools failed: %v", err)
	} else {
		fmt.Printf("\nAvailable tools (%d):\n", len(toolsResult.Tools))
		for _, t := range toolsResult.Tools {
			fmt.Printf("  • %s : %s\n", t.Name, t.Description)
		}
	}

	// 5. Example: call a tool (uncomment/adapt to real tool name & args)
	callCtx, callCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer callCancel()

	result, err := session.CallTool(callCtx, &mcp.CallToolParams{
		Name: "search_product_recalls",
		Arguments: map[string]any{
			"query": "contaminated",
			"limit": 1,
		},
	})
	if err != nil {
		log.Printf("CallTool failed: %v", err)
	} else {
		fmt.Printf("\nTool result:\n")
		for _, content := range result.Content {
			if textContent, ok := content.(*mcp.TextContent); ok {
				fmt.Printf("  %s\n", textContent.Text)
			}
		}
	}
}
