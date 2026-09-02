# rk-mcp

Go client and examples for [Recall Kitchen](https://recallkitchen.com/) MCP.

Recall Kitchen searches U.S. CPSC, FDA food, FDA MedWatch, and USDA product recalls. Vehicle/VIN search is coming soon.

## MCP server

```
https://app.recallkitchen.com/mcp
```

Local: `http://localhost:8080/mcp`. Server card: `https://app.recallkitchen.com/.well-known/mcp/server-card.json`.

Tool reference on the product docs: [recallkitchen.com/docs/#mcp](https://recallkitchen.com/docs/#mcp). Implementation notes in the main repo: [docs/mcp.md](https://github.com/Recall-Kitchen/recall-kitchen/blob/master/docs/mcp.md).

## Auth

Two ways to call tools:

1. **API key (recommended)** — free, higher limits. Create a key in the app under [Integrations](https://app.recallkitchen.com/#/integrations). Send `X-API-Key: rk_...` (required for Grok) or `Authorization: Bearer rk_...` (Claude Code, Cursor).
2. **x402** — anonymous USDC on Base (~$0.025 / call). Set `X402_EVM_PRIVATE_KEY`. Watch, inventory, and notification tools still need an API key.

Grok's HTTP MCP client treats `Authorization: Bearer` as OAuth. Use `X-API-Key` for Grok.

```toml
# ~/.grok/config.toml
[mcp_servers.recall-kitchen]
url = "https://app.recallkitchen.com/mcp"
enabled = true

[mcp_servers.recall-kitchen.headers]
X-API-Key = "rk_..."
```

```bash
export RK_API_KEY=rk_...
```

## Tools

**Public** (API key or x402)

| Tool | Use |
|---|---|
| `search_product_recalls` | Keyword search. Optional `source`, `since`/`until` (`YYYY-MM-DD`), `location`, `offset`, `limit` (1–100, default 3). |
| `search_recalls_by_identifier` | Exact UPC, lot, model, or extracted product name. |
| `search_product_recalls_by_upc` | Look up a UPC and match extracted recall UPCs. `found=false` if unknown. |
| `lookup_product` | Catalog / USDA branded-food info. Does not search recalls. |
| `search_product_recalls_from_image` | Public HTTPS URL, `data:image/...;base64` URI, or MCP image content. No local file paths. |
| `get_product_recall` | Full description and extracted lots, UPCs, models, locations, contacts, image URLs. |

Search results truncate descriptions and cap extracted products at 5. Call `get_product_recall` for the full notice.

`location` is ANDed with the query. `CA` matches California / "northern california". `United States` matches U.S. notices, not Canada/Ontario.

**Account** (API key required)

`list_watch_patterns`, `add_watch_pattern`, `remove_watch_pattern`, `list_inventory`, `add_inventory_product`, `remove_inventory_product`, `check_tracked_products`, `list_recall_notifications`.

Watch patterns use the same syntax as the search box: words are AND, `OR` is or, `-term` excludes, `"quoted phrase"` is a phrase.

Prompts: `check_product`, `check_upc`, `scan_image`. Resources: `recall://docs/tools`, `recall://docs/sources`.

## Go client

```go
cc, err := rkmcp.NewClient(rkmcp.Config{
    ServerURL: "https://app.recallkitchen.com/mcp",
    APIKey:    os.Getenv("RK_API_KEY"),
})
recalls, err := cc.SearchProductRecalls(ctx, "contamination", 3)
detail, err := cc.GetProductRecall(ctx, recalls[0].ID)
hits, err := cc.SearchRecallsByIdentifier(ctx, rkmcp.IdentifierOptions{LotCode: "I31C"})
```

`APIKey` is also read from `RK_API_KEY` or `RECALL_KITCHEN_API_KEY`. The client sends `X-API-Key`. Anonymous clients can set `EVMPrivateKey` / `X402_EVM_PRIVATE_KEY` instead.

Module: `github.com/Recall-Kitchen/rk-mcp/go`.

## Examples

```bash
git clone https://github.com/Recall-Kitchen/rk-mcp.git
cd rk-mcp/examples/go
export RK_API_KEY=rk_...
go run ./search_product_recalls -query "spinach" -location Iowa
go run ./get_product_recall            # searches, then fetches full text
go run ./search_by_identifier -name "infant formula"
```

See [examples/go](examples/go/).

## Other clients

- **Grok**: `X-API-Key` header (see above).
- **Claude Code**: `claude mcp add --transport http recall-kitchen https://app.recallkitchen.com/mcp --header "Authorization: Bearer ${RK_API_KEY}"`
- **Cursor**: `headers.Authorization` = `Bearer ${env:RECALL_KITCHEN_API_KEY}` in `mcp.json`.

## Support

[support@recallkitchen.com](mailto:support@recallkitchen.com)
