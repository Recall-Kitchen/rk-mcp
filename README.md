# MCP Examples

This repository contains examples demonstrating the usage of the Model Context Protocol (MCP) for [Recall Kitchen](https://recallkitchen.com/), including searching recalls.

## MCP Server

Recall Kitchen hosts a MCP server at:

```
https://app.recallkitchen.com/mcp
```

Currently we offer the following tools:

- `search_product_recalls(query: string, limit: int)`

## Getting Started

1. Clone this repository
2. Navigate to the desired example directory
3. Follow the instructions in the example's README

## Project Structure

This repository is organized as follows:

- Go
   - [`examples/go/search_product_recalls/`](examples/go/search_product_recalls/): A client example for calling the `search_product_recalls` tool.

## Contributing

Feel free to contribute new examples by opening a pull request.
