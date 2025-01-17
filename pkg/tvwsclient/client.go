package tvwsclient

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

// Client represents a TradingView WebSocket client
type Client struct {
	ws            *websocket.Conn
	requestHeader http.Header
	wsURL         string
	validateURL   string
	authToken     string
	maxRetries    int
}

var heartbeatRegex = regexp.MustCompile(`~h~\d+`)

// NewClient creates a new TradingView WebSocket client
func NewClient(authToken string, options ...Option) (*Client, error) {
	client := &Client{
		requestHeader: defaultHeaders(),
		wsURL:         "wss://prodata.tradingview.com/socket.io/websocket?from=screener%2F",
		validateURL:   "https://scanner.tradingview.com/symbol?symbol=%s%%3A%s&fields=market&no_404=false",
		authToken:     authToken,
		maxRetries:    5,
	}

	// Apply options
	for _, opt := range options {
		opt(client)
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes a WebSocket connection
func (c *Client) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, c.requestHeader)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	c.ws = conn
	return c.SendInitMessage()
}

// reconnect attempts to reconnect to the WebSocket server
func (c *Client) reconnect() error {
	if c.ws != nil {
		c.ws.Close()
	}

	if err := c.connect(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SendInitMessage() error {
	if err := SendSetAuthTokenMessage(c, c.authToken); err != nil {
		return err
	}

	if err := SendSetLocalMessage(c); err != nil {
		return err
	}

	return nil
}

func (c *Client) ReadMessage(dataChan chan<- TVResponse) error {
	retries := 0
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}

			// Attempt to reconnect if we haven't exceeded max retries
			if retries < c.maxRetries {
				retries++
				slog.Error("connection error, attempting reconnect",
					"attempt", retries,
					"max_retries", c.maxRetries,
					"error", err)

				if err := c.reconnect(); err != nil {
					slog.Error("reconnection attempt failed",
						"attempt", retries,
						"error", err)
					continue
				}

				slog.Info("successfully reconnected",
					"attempt", retries)
				continue
			}

			return fmt.Errorf("error reading message after %d reconnection attempts: %w", retries, err)
		}

		// Reset retry counter on successful message
		retries = 0

		// Handle heartbeat messages
		if heartbeatRegex.Match(message) {
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				slog.Error("error sending heartbeat response",
					"error", err)
				continue // Don't return on heartbeat errors, try to reconnect
			}
			continue
		}

		// Process data messages
		parts := strings.Split(string(message), "~m~")
		for _, part := range parts {
			if strings.HasPrefix(part, "{") {
				slog.Debug("received", "message", part)
				var response TVResponse
				if err := json.Unmarshal([]byte(part), &response); err != nil {
					continue
				}
				dataChan <- response
			}
		}
	}
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
