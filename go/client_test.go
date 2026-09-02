package rkmcp_test

import (
	"cmp"
	"context"
	"os"
	"testing"
	"time"

	rkmcp "github.com/Recall-Kitchen/rk-mcp/go"

	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T) rkmcp.Client {
	t.Helper()
	if testing.Short() {
		t.Skip("-short flag provided")
	}

	apiKey := cmp.Or(os.Getenv("RK_API_KEY"), os.Getenv("RECALL_KITCHEN_API_KEY"))
	cc, err := rkmcp.NewClient(rkmcp.Config{
		ServerURL: cmp.Or(os.Getenv("MCP_SERVER_URL"), "https://app.recallkitchen.com/mcp"),
		Timeout:   20 * time.Second,
		APIKey:    apiKey,
	})
	if apiKey == "" {
		t.Skip("set RK_API_KEY or RECALL_KITCHEN_API_KEY to run live MCP tests")
	}
	require.NoError(t, err)
	t.Cleanup(func() { _ = cc.Close() })
	return cc
}

func TestClient_SearchProductRecalls(t *testing.T) {
	cc := testClient(t)

	recalls, err := cc.SearchProductRecalls(context.Background(), "contamination", 1)
	require.NoError(t, err)
	require.NotEmpty(t, recalls)

	recall := recalls[0]
	require.NotEmpty(t, recall.ID)
	require.NotEmpty(t, recall.Title)
	t.Logf("id=%s source=%s title=%s", recall.ID, recall.Source, recall.Title)
}

func TestClient_GetProductRecall(t *testing.T) {
	cc := testClient(t)

	found, err := cc.SearchProductRecalls(context.Background(), "contamination", 1)
	require.NoError(t, err)
	require.NotEmpty(t, found)

	detail, err := cc.GetProductRecall(context.Background(), found[0].ID)
	require.NoError(t, err)
	require.Equal(t, found[0].ID, detail.ID)
	require.NotEmpty(t, detail.Description)
}

func TestClient_SearchRecallsByIdentifier(t *testing.T) {
	cc := testClient(t)

	found, err := cc.SearchProductRecalls(context.Background(), "chicken", 5)
	require.NoError(t, err)
	require.NotEmpty(t, found)

	opts := rkmcp.IdentifierOptions{Limit: 3}
	for _, hit := range found {
		detail, err := cc.GetProductRecall(context.Background(), hit.ID)
		require.NoError(t, err)
		if detail.Extracted == nil {
			continue
		}
		for _, p := range detail.Extracted.Products {
			if len(p.LotCodes) > 0 {
				opts.LotCode = p.LotCodes[0]
			}
			if len(p.UPCs) > 0 {
				opts.UPC = p.UPCs[0]
			}
			if opts.LotCode != "" || opts.UPC != "" {
				break
			}
		}
		if opts.LotCode != "" || opts.UPC != "" {
			break
		}
	}
	if opts.LotCode == "" && opts.UPC == "" {
		t.Skip("no extracted lot or UPC in search hits")
	}

	res, err := cc.SearchRecallsByIdentifier(context.Background(), opts)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.Recalls)
}
