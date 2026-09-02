# Go examples

These programs call Recall Kitchen MCP with the [`rk-mcp/go`](../../go) client.

Set a free API key (preferred) or an x402 private key:

```bash
export RK_API_KEY=rk_...
# or
export X402_EVM_PRIVATE_KEY=0x...
```

Create a key at [app.recallkitchen.com/#/integrations](https://app.recallkitchen.com/#/integrations).

| Example | What it shows |
|---|---|
| [search_product_recalls](search_product_recalls/) | Keyword search with optional location |
| [get_product_recall](get_product_recall/) | Full notice after a search hit |
| [search_by_identifier](search_by_identifier/) | Lot / UPC / model / product name |

```bash
cd examples/go
go run ./search_product_recalls -query spinach -location Iowa
go run ./get_product_recall
go run ./search_by_identifier -name "infant formula"
```

Tool list and auth notes: [root README](../../README.md). Product docs: [recallkitchen.com/docs/#mcp](https://recallkitchen.com/docs/#mcp).
