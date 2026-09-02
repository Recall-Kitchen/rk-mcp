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
		serverURL   string
		upc         string
		lot         string
		model       string
		productName string
		limit       int
	)
	flag.StringVar(&serverURL, "address", "https://app.recallkitchen.com/mcp", "MCP server HTTP endpoint")
	flag.StringVar(&upc, "upc", "", "UPC or EAN")
	flag.StringVar(&lot, "lot", "", "lot or batch code")
	flag.StringVar(&model, "model", "", "model number")
	flag.StringVar(&productName, "name", "infant formula", "extracted product name")
	flag.IntVar(&limit, "limit", 3, "max results")
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

	res, err := cc.SearchRecallsByIdentifier(context.Background(), rkmcp.IdentifierOptions{
		UPC:         upc,
		LotCode:     lot,
		ModelNumber: model,
		ProductName: productName,
		Limit:       limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: search_recalls_by_identifier: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d recalls\n", len(res.Recalls))
	for _, recall := range res.Recalls {
		fmt.Printf("\nID: %s\nTitle: %s\n", recall.ID, recall.Title)
	}
}
