# MCP Go Examples for Recall Kitchen

This directory contains beginner-friendly Go examples demonstrating how to use the Model Context Protocol (MCP) with [Recall Kitchen](https://recallkitchen.com/). These examples show how to search for product recalls using simple Go code.

## Key Concepts
- **MCP (Model Context Protocol)**: A simple protocol for calling tools over the web, like querying a database. Think of it as a remote function call.
- **Recall Kitchen Tools**: Services like `search_product_recalls` let you find safety recall info for products. Some require payment via x402 (a crypto-based system).
- **Go Client**: The code here uses the `rkmcp` package to connect and query without you having to handle low-level details.

## Examples
### search_product_recalls
- **Location**: `search_product_recalls/`
- **Description**: A simple program that calls the `search_product_recalls` tool to query recalls (e.g., for "bacteria"). It demonstrates JSON-RPC requests and response handling.
- **What You'll Learn**: How to set up a client, make a paid/free query, and print results.

## Getting Started (Detailed Steps)
1. **Install Go**: If you don't have Go, download it from [go.dev](https://go.dev/doc/install) (version 1.21+ recommended). Run `go version` to check.
2. **Sign Up for Recall Kitchen**: Create an account at [recallkitchen.com](https://recallkitchen.com/) for API access.
3. **Set Up Payments (Optional for Paid Tools)**: For `search_product_recalls`, set the `X402_EVM_PRIVATE_KEY` environment variable with your EVM wallet key. See the root README.md for details.
4. **Clone the Repository**: Run `git clone https://github.com/Recall-Kitchen/rk-mcp.git` and navigate to `examples/go/`.
5. **Install Dependencies**: Run `go mod tidy` to download required packages.
6. **Run an Example**: Go to `search_product_recalls/` and follow its README.md. Typically, `go run main.go`.

## Understanding the Code
Each example includes a `main.go` file that:
- Creates a client with configuration.
- Calls a tool like `SearchProductRecalls`.
- Handles and prints the response.
Look at the godoc comments in the code for more explanations.

## Troubleshooting
- **Error: Missing Private Key**: Set `X402_EVM_PRIVATE_KEY` for paid tools.
- **Connection Issues**: Check your internet and the server URL.
- See FAQ.md in the root for more help.

## Commands
- [`search_product_recalls/`](./search_product_recalls/): Run `go run main.go` to search recalls (edit the query in code).