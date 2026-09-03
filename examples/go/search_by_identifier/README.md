# search_by_identifier

Calls `search_recalls_by_identifier` (extracted UPC, lot, model, product name, or VIN). Multiple flags are AND-matched on the same product.

```bash
export RK_API_KEY=rk_...
go run . -name "infant formula"
go run . -upc 012345678905
go run . -lot I31C
go run . -vin 1HGCM82633A004352
```
