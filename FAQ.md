# FAQ

## Do I have to pay?

No. Create a free API key under [Integrations](https://app.recallkitchen.com/#/integrations) and send `X-API-Key` or `Authorization: Bearer`. Anonymous calls without a key pay with x402 (USDC on Base).

## Why does Grok fail with "Auth required, when send initialize request"?

Grok treats `Authorization: Bearer` as MCP OAuth. Recall Kitchen does not implement OAuth. Use `X-API-Key` instead.

## Why is the description truncated?

Search tools return compact results so agents stay inside context. Call `get_product_recall` with the `id` for the full text and full extracted lists.

## Can I pass a local image path?

No. `search_product_recalls_from_image` accepts a public HTTPS URL, a `data:image/...;base64` URI, or MCP image content. Local and loopback URLs are rejected.

## Why didn't my watch pattern fire?

Patterns use search-box syntax, not substring `Contains`. `Generac Generator -Portable` means Generac AND Generator AND NOT Portable. Set a severity at or above your notification minimum. Location of interest filters state-only notices.

## Vehicle / VIN search?

Coming soon. Current sources are CPSC, FDA food, FDA MedWatch, and USDA.

## More help

[recallkitchen.com/docs](https://recallkitchen.com/docs/) · [docs/mcp.md](https://github.com/Recall-Kitchen/recall-kitchen/blob/master/docs/mcp.md) · [support@recallkitchen.com](mailto:support@recallkitchen.com)
