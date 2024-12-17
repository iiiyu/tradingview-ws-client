package tvwsclient

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client represents a TradingView WebSocket client
type Client struct {
	ws            *websocket.Conn
	requestHeader http.Header
	wsURL         string
	validateURL   string
}

// NewClient creates a new TradingView WebSocket client
func NewClient(options ...Option) (*Client, error) {
	client := &Client{
		requestHeader: defaultHeaders(),
		wsURL:         "wss://data.tradingview.com/socket.io/websocket?from=screener%2F",
		validateURL:   "https://scanner.tradingview.com/symbol?symbol=%s%%3A%s&fields=market&no_404=false",
	}

	// Apply options
	for _, opt := range options {
		opt(client)
	}

	conn, _, err := websocket.DefaultDialer.Dial(client.wsURL, client.requestHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	client.ws = conn

	return client, nil
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	if c.ws != nil {
		return c.ws.Close()
	}
	return nil
}

// defaultHeaders returns the default HTTP headers for the WebSocket connection
func defaultHeaders() http.Header {
	return http.Header{
		"Accept-Encoding": {"gzip, deflate, br, zstd"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Cache-Control":   {"no-cache"},
		"Origin":          {"https://www.tradingview.com"},
		"Pragma":          {"no-cache"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
	}
}
