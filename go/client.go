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

	x402 "github.com/coinbase/x402/go"
	x402mcp "github.com/coinbase/x402/go/mcp"
	evm "github.com/coinbase/x402/go/mechanisms/evm/exact/client"
	evmsigners "github.com/coinbase/x402/go/signers/evm"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client interface {
	Close() error

	SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error)
}

type Config struct {
	ServerURL string
	Timeout   time.Duration

	EVMPrivateKey string
}

var (
	defaultServerURL = "https://app.recallkitchen.com/mcp"
	defaultTimeout   = 15 * time.Second
)

func NewClient(config Config) (Client, error) {
	hostname, _ := os.Hostname()

	clientImpl := &mcpsdk.Implementation{
		Name:    "rk-mcp-go",
		Version: hostname,
	}
	cc := mcpsdk.NewClient(clientImpl, nil)

	// Set defaults if user didn't provide
	config.ServerURL = cmp.Or(config.ServerURL, defaultServerURL)
	config.Timeout = cmp.Or(config.Timeout, defaultTimeout)

	ctx, _ := context.WithTimeout(context.Background(), config.Timeout)

	session, err := cc.Connect(ctx, &mcpsdk.StreamableClientTransport{Endpoint: config.ServerURL}, nil)
	if err != nil {
		return nil, fmt.Errorf("creating MCP client: %w", err)
	}

	out := &client{
		config:  config,
		client:  cc,
		session: session,
	}

	paymentClient, err := createX402PaymentClient(config)
	if err != nil {
		return nil, fmt.Errorf("creating x402 payment client wrapper: %w", err)
	}
	if paymentClient == nil {
		return out, ErrX402NotConfigured
	}

	out.x402Session = x402mcp.NewX402MCPClient(session, paymentClient, x402mcp.Options{
		AutoPayment: x402mcp.BoolPtr(true),
		OnPaymentRequested: func(context x402mcp.PaymentRequiredContext) (bool, error) {
			price := context.PaymentRequired.Accepts[0]
			fmt.Printf("\n💰 Payment required for tool: %s\n", context.ToolName)
			fmt.Printf("   Amount: %s (%s)\n", price.Amount, price.Asset)
			fmt.Printf("   Network: %s\n", price.Network)
			fmt.Printf("   Approving payment...\n")
			return true, nil
		},
	})

	return out, nil
}

var (
	ErrX402NotConfigured = errors.New("x402 client was not configured")
)

func createX402PaymentClient(config Config) (*x402.X402Client, error) {
	privateKey := strings.TrimSpace(cmp.Or(config.EVMPrivateKey, os.Getenv("X402_EVM_PRIVATE_KEY")))
	if privateKey == "" {
		return nil, nil
	}

	evmSigner, err := evmsigners.NewClientSignerFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("problem setting up client signer: %w", err)
	}

	paymentClient := x402.Newx402Client()
	paymentClient.Register("eip155:8453", evm.NewExactEvmScheme(evmSigner, nil)) // mainnet Base

	return paymentClient, nil
}

type client struct {
	config      Config
	client      *mcpsdk.Client
	session     *mcpsdk.ClientSession
	x402Session *x402mcp.X402MCPClient
}

func (c *client) Close() error {
	if c.session != nil {
		c.session.Close()
	}

	return nil
}

func (c *client) SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error) {
	callCtx, callCancel := context.WithTimeout(ctx, 10*time.Second)
	defer callCancel()

	limit = min(cmp.Or(limit, 3), 100) // default 3, but (0,100]
	if limit <= 0 {
		limit = 3
	}

	if c.x402Session != nil {
		return c.x402SearchProductRecalls(callCtx, query, limit)
	}
	return c.mcpSearchProductRecalls(callCtx, query, limit)
}

func (c *client) x402SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error) {
	result, err := c.x402Session.CallTool(ctx, "search_product_recalls", map[string]interface{}{
		"query": query,
		"limit": limit,
	})
	if err != nil {
		return nil, fmt.Errorf("calling tool x402+MCP search_product_recalls: %w", err)
	}

	if result.IsError {
		fmt.Printf("IsError=%v\n", result.IsError)
		for _, content := range result.Content {
			fmt.Printf("\n%v\n", content.Text)
		}

		fmt.Printf("PaymentResponse: %#v\n", result.PaymentResponse)
		fmt.Printf("PaymentMade=%v\n", result.PaymentMade)

		return nil, errors.New("x402+MCP tool call retrned IsError=true")
	}

	if result.PaymentResponse != nil {
		fmt.Printf("\n%#v\n", result.PaymentResponse)

		fmt.Println("\n📦 Payment Receipt:")
		fmt.Printf("   Success: %v\n", result.PaymentResponse.Success)
		if result.PaymentResponse.Transaction != "" {
			fmt.Printf("   Transaction: %s\n", result.PaymentResponse.Transaction)
		}
	}

	for _, content := range result.Content {
		var recalls Recalls
		err := json.NewDecoder(strings.NewReader(content.Text)).Decode(&recalls)
		if err == nil {
			return recalls.Recalls, nil
		}
		return nil, fmt.Errorf("reading x402+MCP search_product_recalls response json: %w", err)
	}

	return nil, errors.New("no response from x402+MCP search_product_recalls found")
}

func (c *client) mcpSearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error) {
	result, err := c.session.CallTool(ctx, &mcpsdk.CallToolParams{
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
		if textContent, ok := content.(*mcpsdk.TextContent); ok {
			var recalls Recalls
			err := json.NewDecoder(strings.NewReader(textContent.Text)).Decode(&recalls)
			if err == nil {
				return recalls.Recalls, nil
			}
			return nil, fmt.Errorf("reading search_product_recalls response json: %w", err)
		}
	}

	return nil, errors.New("no response from MCP search_product_recalls found")
}
