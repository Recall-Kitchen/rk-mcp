# Search Product Recalls MCP Client Example

This example provides a beginner-friendly Go program to call the `search_product_recalls` tool on Recall Kitchen's MCP server using JSON-RPC over HTTP. It's a great starting point to learn how to integrate recall searches into your apps.

## What It Does
The program creates a client, sends a query (e.g., "bacteria") to the server, and prints matching product recalls. The server returns details like title, description, and date in JSON format. This tool requires payment via x402 for use.

## Running the Example
1. **Prerequisites**: Ensure you have Go installed, signed up at [Recall Kitchen](https://recallkitchen.com), and set `X402_EVM_PRIVATE_KEY` for payments (see root README.md).
2. **Install Dependencies**: Run `go mod tidy`.
3. **Run the Program**:
   ```
   go run main.go
   ```
4. The output will display recall details. If payment is needed, it will prompt and handle it.

## Sample Output
Here's what you might see (shortened for example):
```
ID: recall-123
Source: FDA
Title: Peanut Butter Recall Due to Bacteria
PublishedOn: 2023-05-10
Description: Contaminated with salmonella bacteria...
```

## Code Overview
- **main.go**: Parses flags, creates a client using `rkmcp.NewClient`, calls `SearchProductRecalls`, and prints results.
- Key Parts:
  - Configures server URL and timeout.
  - Handles payment if required.
  - Parses JSON response into Recall structs.

See godoc comments in the code for more details.

## Customizing the Query
Edit `main.go` to change the search:
```go
recalls, err := cc.SearchProductRecalls(context.Background(), "your_query_here", 5) // Change query and limit (1-100)
```
- **Query Examples**: "contaminated peanut butter", "toy safety", "vehicle recall".
- **Limit**: Sets max results (default 3).

## Error Handling
Common issues and fixes:
- **Missing Private Key**: Error "missing EVM private key" – Set `X402_EVM_PRIVATE_KEY` env var.
- **Connection Failed**: Check internet, server URL, or timeout. Try increasing timeout in Config.
- **Payment Error**: Ensure your wallet has funds; check console for payment details.
- **No Results**: Try a broader query. If "no response found", verify server status.
- In code, errors are printed; add your own handling like logging.

## Support
Contact us at [support@recallkitchen.com](mailto:support@recallkitchen.com) for help or questions.