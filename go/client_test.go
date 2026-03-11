package rkmcp_test

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	rkmcp "github.com/Recall-Kitchen/rk-mcp/go"

	"github.com/stretchr/testify/require"
)

func TestClient_SearchProductRecalls(t *testing.T) {
	if testing.Short() {
		t.Skip("-short flag provided")
	}

	cc, err := rkmcp.NewClient(rkmcp.Config{
		ServerURL: cmp.Or(os.Getenv("MCP_SERVER_URL"), "https://app.recallkitchen.com/mcp"),
		Timeout:   10 * time.Second,
	})

	if strings.EqualFold(os.Getenv("GITHUB_ACTIONS"), "1") {
		if err != nil && strings.Contains(err.Error(), "x402 client was not configured") {
			t.Skip("full x402 config not setup")
		}
	} else {
		require.NoError(t, err)
	}

	recalls, err := cc.SearchProductRecalls(context.Background(), "bacteria", 1)
	require.NoError(t, err)
	require.Len(t, recalls, 1)

	recall := recalls[0]
	fmt.Printf("\n")
	fmt.Printf("ID: %s\n", recall.ID)
	fmt.Printf("Source: %s\n", recall.Source)
	fmt.Printf("Title: %s\n", recall.Title)
	fmt.Printf("PublishedOn: %s\n", recall.PublishedOn)
	fmt.Printf("\n")
	fmt.Printf("Description: %s\n", recall.Description)
	fmt.Printf("\n")

	require.Contains(t, recall.Description, "bacteria")
}
