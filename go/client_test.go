package rkmcp_test

import (
	"context"
	"fmt"
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
		ServerURL: "http://localhost:8080/mcp",
		Timeout:   10 * time.Second,
	})
	require.NoError(t, err)

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
