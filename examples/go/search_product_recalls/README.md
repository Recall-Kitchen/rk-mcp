# Search Product Recalls MCP Client

This example demonstrates a Go client for calling the `search_product_recalls` tool on an MCP server over HTTP using JSON-RPC.

## What It Does

The client sends a JSON-RPC request to invoke the `search_product_recalls` tool with a query for "apples". The server responds with product recall information related to the query, formatted according to the MCP specification.

## Running the Example

1. Start your MCP server on `https://app.recallkitchen.com/mcp`.

2. Run the client:

   ```bash
   go run main.go
   ```

3. The output will be the tool's response, including recall details in JSON format.

## Code Overview

- Constructs a JSON-RPC 2.0 request with `method: "tools/call"`
- Specifies the tool name and arguments
- Sends POST request to the MCP endpoint
- Parses and displays the response

## Customizing the Query

Edit the `Arguments` map in `main.go` to change the search query:

```go
Arguments: map[string]any{
    "query": "your_query_here",
},
```

## Server Compatibility

This client works with MCP servers that support HTTP-based JSON-RPC for tool calls. Ensure your server implements the `tools/call` method for `search_product_recalls`.
