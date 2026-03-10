package rkmcp

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client interface {
	Close() error

	SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error)
}

type Config struct {
	ServerURL string
	Timeout   time.Duration
}

var (
	defaultServerURL = "https://app.recallkitchen.com/mcp"
	defaultTimeout   = 15 * time.Second
)

func NewClient(config Config) (Client, error) {
	hostname, _ := os.Hostname()

	clientImpl := &mcp.Implementation{
		Name:    "rk-mcp-go",
		Version: hostname,
	}
	cc := mcp.NewClient(clientImpl, nil)

	// Set defaults if user didn't provide
	config.ServerURL = cmp.Or(config.ServerURL, defaultServerURL)
	config.Timeout = cmp.Or(config.Timeout, defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), config.Timeout)

	session, err := cc.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: config.ServerURL}, nil)
	if err != nil {
		return nil, fmt.Errorf("creating MCP client: %w", err)
	}

	return &client{
		config:  config,
		client:  cc,
		session: session,
	}, nil
}

type client struct {
	config  Config
	client  *mcp.Client
	session *mcp.ClientSession
}

func (c *client) Close() error {
	if c.session != nil {
		c.session.Close()
	}

	return nil
}

func (c *client) SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error) {
	callCtx, callCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer callCancel()

	limit = cmp.Or(limit, 3)

	result, err := c.session.CallTool(callCtx, &mcp.CallToolParams{
		Name: "search_product_recalls",
		Arguments: map[string]any{
			"query": query,
			"limit": limit,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("calling tool search_product_recalls: %w", err)
	}

	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			var recalls Recalls
			err := json.NewDecoder(strings.NewReader(textContent.Text)).Decode(&recalls)
			if err == nil {
				return recalls.Recalls, nil
			}
			return nil, fmt.Errorf("reading search_product_recalls response json: %w", err)
		}
	}

	return nil, errors.New("no response from search_product_recalls found")
}
