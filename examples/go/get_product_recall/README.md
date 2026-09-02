# get_product_recall

Fetches one recall by id, including extracted lots, UPCs, locations, and contact info.

```bash
export RK_API_KEY=rk_...
go run . -id RECALL_ID
# or search first, then fetch:
go run . -query contamination
```
