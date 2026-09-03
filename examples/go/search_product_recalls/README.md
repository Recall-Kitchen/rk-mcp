# search_product_recalls

Calls `search_product_recalls`. Descriptions are truncated; use [get_product_recall](../get_product_recall/) for the full text.

```bash
export RK_API_KEY=rk_...
go run . -query spinach -location Iowa -limit 3
```

Flags: `-address`, `-query`, `-source` (`nhtsa`, `cpsc`, ...), `-location`, `-limit`.

Auth: `RK_API_KEY` / `RECALL_KITCHEN_API_KEY`, or `X402_EVM_PRIVATE_KEY` for anonymous USDC on Base. Without either, the program prints setup help and exits 0 (so CI can compile-run it).
