package rkmcp

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	x402 "github.com/x402-foundation/x402/go"
	x402mcp "github.com/x402-foundation/x402/go/mcp"
	evm "github.com/x402-foundation/x402/go/mechanisms/evm/exact/client"
	evmsigners "github.com/x402-foundation/x402/go/signers/evm"
)

type Client interface {
	Close() error

	SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error)
	SearchProductRecallsOpts(ctx context.Context, opts SearchOptions) (*SearchResult, error)
	GetProductRecall(ctx context.Context, recallID string) (*Recall, error)
	SearchRecallsByIdentifier(ctx context.Context, opts IdentifierOptions) (*SearchResult, error)
}

type Config struct {
	ServerURL string
	Timeout   time.Duration

	// APIKey is a Recall Kitchen rk_ key. Sent as X-API-Key (Grok-safe; also
	// accepted as Authorization: Bearer). When set, search tools skip x402.
	APIKey string

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

	config.ServerURL = cmp.Or(config.ServerURL, defaultServerURL)
	config.Timeout = cmp.Or(config.Timeout, defaultTimeout)
	config.APIKey = strings.TrimSpace(cmp.Or(config.APIKey, os.Getenv("RK_API_KEY"), os.Getenv("RECALL_KITCHEN_API_KEY")))

	transport := &mcpsdk.StreamableClientTransport{
		Endpoint:             config.ServerURL,
		DisableStandaloneSSE: true,
	}
	if config.APIKey != "" {
		transport.HTTPClient = &http.Client{
			Timeout:   config.Timeout,
			Transport: apiKeyRoundTripper{key: config.APIKey, base: http.DefaultTransport},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	session, err := cc.Connect(ctx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("creating MCP client: %w", err)
	}

	out := &client{
		config:  config,
		client:  cc,
		session: session,
	}

	if config.APIKey != "" {
		return out, nil
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
			fmt.Printf("\nPayment required for tool: %s\n", context.ToolName)
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

type apiKeyRoundTripper struct {
	key  string
	base http.RoundTripper
}

func (t apiKeyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Set("X-API-Key", t.key)
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(r)
}

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

func clampLimit(limit int) int {
	limit = min(cmp.Or(limit, 3), 100)
	if limit <= 0 {
		return 3
	}
	return limit
}

func (c *client) SearchProductRecalls(ctx context.Context, query string, limit int) ([]Recall, error) {
	res, err := c.SearchProductRecallsOpts(ctx, SearchOptions{Query: query, Limit: limit})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.Recalls, nil
}

func (c *client) SearchProductRecallsOpts(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	args := map[string]any{}
	if q := strings.TrimSpace(opts.Query); q != "" {
		args["query"] = q
	}
	if s := strings.TrimSpace(opts.Source); s != "" {
		args["source"] = s
	}
	if s := strings.TrimSpace(opts.Since); s != "" {
		args["since"] = s
	}
	if s := strings.TrimSpace(opts.Until); s != "" {
		args["until"] = s
	}
	if s := strings.TrimSpace(opts.Location); s != "" {
		args["location"] = s
	}
	if opts.Offset > 0 {
		args["offset"] = opts.Offset
	}
	args["limit"] = clampLimit(opts.Limit)

	var out SearchResult
	if err := c.callTool(ctx, "search_product_recalls", args, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) GetProductRecall(ctx context.Context, recallID string) (*Recall, error) {
	recallID = strings.TrimSpace(recallID)
	if recallID == "" {
		return nil, fmt.Errorf("recall_id is required")
	}
	var out Recall
	if err := c.callTool(ctx, "get_product_recall", map[string]any{"recall_id": recallID}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) SearchRecallsByIdentifier(ctx context.Context, opts IdentifierOptions) (*SearchResult, error) {
	args := map[string]any{}
	if s := strings.TrimSpace(opts.UPC); s != "" {
		args["upc"] = s
	}
	if s := strings.TrimSpace(opts.LotCode); s != "" {
		args["lot_code"] = s
	}
	if s := strings.TrimSpace(opts.ModelNumber); s != "" {
		args["model_number"] = s
	}
	if s := strings.TrimSpace(opts.ProductName); s != "" {
		args["product_name"] = s
	}
	if s := strings.TrimSpace(opts.VIN); s != "" {
		args["vin"] = s
	}
	if opts.Offset > 0 {
		args["offset"] = opts.Offset
	}
	args["limit"] = clampLimit(opts.Limit)
	if len(args) == 1 { // only limit
		return nil, fmt.Errorf("upc, lot_code, model_number, product_name, or vin is required")
	}

	var out SearchResult
	if err := c.callTool(ctx, "search_recalls_by_identifier", args, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) callTool(ctx context.Context, name string, args map[string]any, dest any) error {
	callCtx, callCancel := context.WithTimeout(ctx, cmp.Or(c.config.Timeout, 10*time.Second))
	defer callCancel()

	if c.x402Session != nil {
		result, err := c.x402Session.CallTool(callCtx, name, args)
		if err != nil {
			return fmt.Errorf("calling tool x402+MCP %s: %w", name, err)
		}
		if result.IsError {
			var bits []string
			for _, content := range result.Content {
				bits = append(bits, content.Text)
			}
			return fmt.Errorf("x402+MCP %s: %s", name, strings.Join(bits, " "))
		}
		for _, content := range result.Content {
			if err := json.NewDecoder(strings.NewReader(content.Text)).Decode(dest); err != nil {
				return fmt.Errorf("reading x402+MCP %s response json: %w", name, err)
			}
			return nil
		}
		return fmt.Errorf("no response from x402+MCP %s", name)
	}

	result, err := c.session.CallTool(callCtx, &mcpsdk.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		return fmt.Errorf("calling tool %s: %w", name, err)
	}
	if result.IsError {
		var bits []string
		for _, content := range result.Content {
			if textContent, ok := content.(*mcpsdk.TextContent); ok {
				bits = append(bits, textContent.Text)
			}
		}
		return fmt.Errorf("%s: %s", name, strings.Join(bits, " "))
	}
	for _, content := range result.Content {
		if textContent, ok := content.(*mcpsdk.TextContent); ok {
			if err := json.NewDecoder(strings.NewReader(textContent.Text)).Decode(dest); err != nil {
				return fmt.Errorf("reading %s response json: %w", name, err)
			}
			return nil
		}
	}
	return fmt.Errorf("no response from MCP %s", name)
}
