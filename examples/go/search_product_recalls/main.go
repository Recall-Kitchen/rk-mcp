package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	rkmcp "github.com/Recall-Kitchen/rk-mcp/go"
)

func main() {
	var (
		serverURL string
		query     string
		location  string
		limit     int
	)
	flag.StringVar(&serverURL, "address", "https://app.recallkitchen.com/mcp", "MCP server HTTP endpoint")
	flag.StringVar(&query, "query", "contamination", "search query")
	flag.StringVar(&location, "location", "", "optional location filter (e.g. Iowa)")
	flag.IntVar(&limit, "limit", 3, "max results (1-100)")
	flag.Parse()

	cc, err := rkmcp.NewClient(rkmcp.Config{
		ServerURL: cmp.Or(os.Getenv("MCP_SERVER_URL"), serverURL),
		Timeout:   20 * time.Second,
		// APIKey from RK_API_KEY / RECALL_KITCHEN_API_KEY, or:
		// EVMPrivateKey / X402_EVM_PRIVATE_KEY for anonymous x402.
	})
	if err != nil && errors.Is(err, rkmcp.ErrX402NotConfigured) {
		fmt.Fprintln(os.Stderr, "Set RK_API_KEY (free) or X402_EVM_PRIVATE_KEY (anonymous USDC on Base).")
		fmt.Fprintln(os.Stderr, "Create a key at https://app.recallkitchen.com/#/integrations")
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	res, err := cc.SearchProductRecallsOpts(context.Background(), rkmcp.SearchOptions{
		Query:    query,
		Location: location,
		Limit:    limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: searching product recalls: %v\n", err)
		os.Exit(1)
	}

	for _, recall := range res.Recalls {
		fmt.Printf("ID: %s\n", recall.ID)
		fmt.Printf("Source: %s\n", recall.Source)
		fmt.Printf("Title: %s\n", recall.Title)
		fmt.Printf("PublishedOn: %s\n", recall.PublishedOn.Format(time.DateOnly))
		fmt.Printf("URL: %s\n", recall.URL)
		if recall.Extracted != nil && len(recall.Extracted.Locations) > 0 {
			fmt.Printf("Locations: %v\n", recall.Extracted.Locations)
		}
		fmt.Printf("\n%s\n\n", recall.Description)
	}
}
