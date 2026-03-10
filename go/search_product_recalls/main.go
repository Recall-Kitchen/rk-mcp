package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Request struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

type Params struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      int         `json:"id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	// Create the request as per the curl
	req := Request{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: Params{
			Name: "search_product_recalls",
			Arguments: map[string]any{
				"query": "apples",
			},
		},
		ID: 1,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	// Send POST request
	address := cmp.Or(os.Getenv("MCP_SERVER_ADDRESS"), "https://app.recallkitchen.com/mcp")
	resp, err := http.Post(address, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Parse response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		log.Fatalf("RPC error: %s", response.Error.Message)
	}

	// Print the result
	resultJSON, err := json.MarshalIndent(response.Result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(resultJSON))
}
