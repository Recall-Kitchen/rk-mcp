# Contributing

1. Fork and branch from `master`.
2. Keep the Go client backward compatible (`SearchProductRecalls(ctx, query, limit)` stays).
3. `gofmt` the `go/` package. Run `cd go && go test ./...` (live tests skip without `RK_API_KEY`).
4. `cd examples/go && go run ./search_product_recalls` should print the API-key hint and exit 0 when `RK_API_KEY` is unset.
5. Open a PR against `master`.
