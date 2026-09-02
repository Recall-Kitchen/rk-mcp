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
		recallID  string
		query     string
	)
	flag.StringVar(&serverURL, "address", "https://app.recallkitchen.com/mcp", "MCP server HTTP endpoint")
	flag.StringVar(&recallID, "id", "", "recall id (if empty, search -query first)")
	flag.StringVar(&query, "query", "contamination", "search query used when -id is empty")
	flag.Parse()

	cc, err := rkmcp.NewClient(rkmcp.Config{
		ServerURL: cmp.Or(os.Getenv("MCP_SERVER_URL"), serverURL),
		Timeout:   20 * time.Second,
	})
	if err != nil && errors.Is(err, rkmcp.ErrX402NotConfigured) {
		fmt.Fprintln(os.Stderr, "Set RK_API_KEY or X402_EVM_PRIVATE_KEY.")
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer cc.Close()

	if recallID == "" {
		hits, err := cc.SearchProductRecalls(context.Background(), query, 1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: search: %v\n", err)
			os.Exit(1)
		}
		if len(hits) == 0 {
			fmt.Fprintln(os.Stderr, "no search hits")
			os.Exit(1)
		}
		recallID = hits[0].ID
		fmt.Printf("Using first search hit %s\n\n", recallID)
	}

	recall, err := cc.GetProductRecall(context.Background(), recallID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: get_product_recall: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ID: %s\n", recall.ID)
	fmt.Printf("Source: %s\n", recall.Source)
	fmt.Printf("Title: %s\n", recall.Title)
	fmt.Printf("PublishedOn: %s\n", recall.PublishedOn.Format(time.RFC3339))
	fmt.Printf("URL: %s\n", recall.URL)
	if recall.Extracted != nil {
		if len(recall.Extracted.Locations) > 0 {
			fmt.Printf("Locations: %v\n", recall.Extracted.Locations)
		}
		for _, p := range recall.Extracted.Products {
			fmt.Printf("Product: %s lots=%v upcs=%v\n", p.Name, p.LotCodes, p.UPCs)
		}
	}
	fmt.Printf("\n%s\n", recall.Description)
}
