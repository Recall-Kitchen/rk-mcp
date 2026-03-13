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
	var serverURL string
	var timeout time.Duration

	flag.StringVar(&serverURL, "address", "https://app.recallkitchen.com/mcp", "MCP server HTTP endpoint")
	flag.DurationVar(&timeout, "timeout", 15*time.Second, "request timeout")
	flag.Parse()

	cc, err := rkmcp.NewClient(rkmcp.Config{
		ServerURL:     cmp.Or(os.Getenv("MCP_SERVER_URL"), serverURL),
		Timeout:       10 * time.Second,
		EVMPrivateKey: "", // or set X402_EVM_PRIVATE_KEY env var
	})
	if err != nil && errors.Is(err, rkmcp.ErrX402NotConfigured) {
		fmt.Printf("ERROR: missing EVM private key (Config.EVMPrivateKey or X402_EVM_PRIVATE_KEY in the environment)")
		os.Exit(1)
	}

	recalls, err := cc.SearchProductRecalls(context.Background(), "bacteria", 1)
	if err != nil {
		fmt.Printf("ERROR: searching product recalls: %v\n", err)
		os.Exit(1)
	}

	for _, recall := range recalls {
		fmt.Printf("\n")
		fmt.Printf("ID: %s\n", recall.ID)
		fmt.Printf("Source: %s\n", recall.Source)
		fmt.Printf("Title: %s\n", recall.Title)
		fmt.Printf("PublishedOn: %s\n", recall.PublishedOn)
		fmt.Printf("\n")
		fmt.Printf("Description: %s\n", recall.Description)
		fmt.Printf("\n")
	}
}
