package tvwsclient

import (
	"encoding/json"
	"fmt"
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
}

var heartbeatRegex = regexp.MustCompile(`~h~\d+`)

// NewClient creates a new TradingView WebSocket client
func NewClient(authToken string, options ...Option) (*Client, error) {
	client := &Client{
		requestHeader: defaultHeaders(),
		wsURL:         "wss://prodata.tradingview.com/socket.io/websocket?from=screener%2F",
		validateURL:   "https://scanner.tradingview.com/symbol?symbol=%s%%3A%s&fields=market&no_404=false",
		authToken:     authToken,
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

func (c *Client) SendInitMessage() error {
	if err := SendSetAuthTokenMessage(c, c.authToken); err != nil {
		return err
	}

	if err := SendSetLocalMessage(c); err != nil {
		return err
	}

	csSession := GenerateSession("cs_")
	if err := SendChartCreateSessionMessage(c, csSession); err != nil {
		return err
	}
	return nil
}

func (c *Client) ReadMessage(dataChan chan<- map[string]interface{}) error {
	// Read messages in a loop
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
			return fmt.Errorf("error reading message: %w", err)
		}

		// Handle heartbeat messages
		if heartbeatRegex.Match(message) {
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return fmt.Errorf("error sending heartbeat response: %w", err)
			}
			continue
		}

		// Process data messages
		parts := strings.Split(string(message), "~m~")
		for _, part := range parts {
			if strings.HasPrefix(part, "{") {
				var response TVResponse
				if err := json.Unmarshal([]byte(part), &response); err != nil {
					continue
				}

				// Only process quote data messages
				if response.Method == "qsd" && len(response.Params) >= 2 {
					// Extract the quote data from params
					if quoteDataRaw, err := json.Marshal(response.Params[1]); err == nil {
						var quote QuoteData
						if err := json.Unmarshal(quoteDataRaw, &quote); err == nil {
							// Convert to map for compatibility with existing channel
							dataMap := map[string]interface{}{
								"m": response.Method,
								"p": []interface{}{response.Params[0], quote},
							}
							dataChan <- dataMap
						}
					}
				}
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
