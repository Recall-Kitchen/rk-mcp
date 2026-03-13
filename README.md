# MCP Examples for Recall Kitchen

Welcome! This repository provides beginner-friendly examples for using the Model Context Protocol (MCP) with [Recall Kitchen](https://recallkitchen.com/), a service that helps you search and access product recall information easily.

## What is Recall Kitchen?
Recall Kitchen is an online platform that aggregates and provides access to product recall data from various sources. It's useful for developers building apps that need to check for safety recalls on products like food, vehicles, or consumer goods.

## What is MCP?
MCP stands for Model Context Protocol. It's a simple way for AI models and services to communicate and call tools over the web, similar to JSON-RPC. Recall Kitchen uses MCP to let you query their recall database securely.

## Prerequisites
Before getting started:
1. **Install Go**: Download and install Go (version 1.21 or later) from [go.dev](https://go.dev/doc/install). Verify with `go version`.
2. **Sign up for Recall Kitchen**: Create a free account at [recallkitchen.com](https://recallkitchen.com/) to get API access.
3. **Set up x402 Payments (for paid tools)**: If using paid features, you'll need an EVM-compatible wallet (e.g., MetaMask) and set the `X402_EVM_PRIVATE_KEY` environment variable. Learn more about x402 [here](https://x402.org/).

## MCP Server
Recall Kitchen hosts an MCP server at:
```
https://app.recallkitchen.com/mcp
```

Currently, we offer the following tools:

**Payment Required**
Recall Kitchen offers MCP tools that require payment via x402 (a secure payment protocol):
- `search_product_recalls(query: string, limit: int)`: Searches for product recalls matching your query (e.g., \"contaminated peanut butter\"). Limit sets the max results.

**Free Tools**
Recall Kitchen offers some tools for free:
- TODO (Check the Recall Kitchen docs for updates).

## Quick Start
1. Clone this repository: `git clone https://github.com/Recall-Kitchen/rk-mcp.git`
2. Navigate to the root: `cd rk-mcp`
3. Install dependencies: `go mod tidy` (in the `go/` or `examples/go/` directories as needed)
4. Run an example: Go to `examples/go/search_product_recalls/` and follow its README.
5. Example command: `go run main.go` (after setting up your environment).

## Getting Started (Detailed)
1. Clone this repository as shown above.
2. Navigate to the desired example directory (e.g., `examples/go/search_product_recalls/`).
3. Follow the instructions in the example's README.md file.
4. If you encounter errors, check the FAQ.md (in root) or contact support.

## Project Structure
This repository is organized as follows:
- **go/**: Core client library for MCP interactions.
- **examples/go/search_product_recalls/**: A simple client example for calling the `search_product_recalls` tool.

## Contributing
Feel free to contribute new examples or improvements! See CONTRIBUTING.md for guidelines. We welcome beginners – start by forking the repo and submitting a pull request.
