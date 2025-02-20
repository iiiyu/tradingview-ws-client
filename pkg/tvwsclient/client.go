package tvwsclient

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a TradingView WebSocket client
type Client struct {
	ws            *websocket.Conn
	mu            sync.Mutex // protects ws
	requestHeader http.Header
	wsURL         string
	maxRetries    int
	done          chan struct{} // Channel to signal connection close
	reconnecting  bool          // Flag to indicate reconnection in progress
	pingInterval  time.Duration // Interval for sending ping messages
	writeTimeout  time.Duration // Timeout for write operations
	readTimeout   time.Duration // Timeout for read operations
	isConnected   bool
}

var heartbeatRegex = regexp.MustCompile(`~h~\d+`)

// NewClient creates a new TradingView WebSocket client
func NewClient(options ...Option) (*Client, error) {
	client := &Client{
		requestHeader: defaultHeaders(),
		wsURL:         "wss://prodata.tradingview.com/socket.io/websocket?from=screener%2F",
		maxRetries:    5,
		done:          make(chan struct{}),
		pingInterval:  30 * time.Second,
		writeTimeout:  60 * time.Second,
		readTimeout:   60 * time.Second,
	}

	// Apply options
	for _, opt := range options {
		opt(client)
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	// Start ping handler
	go client.pingHandler()

	return client, nil
}

// connect establishes a WebSocket connection
func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ws != nil {
		c.ws.Close()
		c.ws = nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(c.wsURL, c.requestHeader)
	if err != nil {
		c.isConnected = false
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.ws = conn
	c.isConnected = true

	// Setup ping handler to respond to server pings
	c.ws.SetPingHandler(func(appData string) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		if !c.isConnected || c.ws == nil {
			return fmt.Errorf("connection is closed")
		}
		return c.ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
	})

	return c.SendInitMessage()
}

// pingHandler sends periodic ping messages to keep the connection alive
func (c *Client) pingHandler() {
	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			if !c.isConnected || c.ws == nil {
				c.mu.Unlock()
				continue
			}
			err := c.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(c.writeTimeout))
			c.mu.Unlock()

			if err != nil {
				slog.Error("ping error", "error", err)
				if err := c.reconnect(); err != nil {
					slog.Error("failed to reconnect after ping error", "error", err)
				}
			}
		case <-c.done:
			return
		}
	}
}

// reconnect attempts to reconnect to the WebSocket server
func (c *Client) reconnect() error {
	c.mu.Lock()
	if c.reconnecting {
		c.mu.Unlock()
		return nil // Already reconnecting
	}
	c.reconnecting = true
	c.isConnected = false
	if c.ws != nil {
		c.ws.Close()
		c.ws = nil
	}
	c.mu.Unlock()

	// Add small delay before reconnecting
	time.Sleep(time.Second)

	err := c.connect()

	c.mu.Lock()
	c.reconnecting = false
	c.mu.Unlock()

	return err
}

func (c *Client) Reconnect() error {
	return c.reconnect()
}

func (c *Client) SendInitMessage() error {
	// Should be initialized AuthTokenManager first
	authToken := GetAuthTokenManager().GetToken()

	if err := SendSetAuthTokenMessage(c, authToken); err != nil {
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
		c.mu.Lock()
		if !c.isConnected || c.ws == nil {
			c.mu.Unlock()
			if retries < c.maxRetries {
				retries++
				slog.Info("connection lost, attempting reconnect",
					"attempt", retries,
					"max_retries", c.maxRetries)

				if err := c.reconnect(); err != nil {
					slog.Error("reconnection attempt failed",
						"attempt", retries,
						"error", err)
					time.Sleep(time.Duration(retries) * time.Second)
					continue
				}
			} else {
				return fmt.Errorf("connection is closed and max retries exceeded")
			}
			continue
		}
		ws := c.ws // Store local copy of ws
		c.mu.Unlock()

		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				slog.Error("connection closed", "error", err)
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
					time.Sleep(time.Duration(retries) * time.Second)
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
			c.mu.Lock()
			if !c.isConnected || c.ws == nil {
				c.mu.Unlock()
				continue
			}
			err := c.ws.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

			if err != nil {
				slog.Error("error sending heartbeat response", "error", err)
				if err := c.reconnect(); err != nil {
					slog.Error("reconnection failed after heartbeat error", "error", err)
				}
				continue
			}
			continue
		}

		// Process data messages
		parts := strings.Split(string(message), "~m~")
		for _, part := range parts {
			if strings.HasPrefix(part, "{") {
				var response TVResponse
				if err := json.Unmarshal([]byte(part), &response); err != nil {
					slog.Error("failed to unmarshal message", "error", err)
					continue
				}
				select {
				case dataChan <- response:
				case <-c.done:
					slog.Error("dataChan closed", "response", response)
					return nil
				}
			}
		}
	}
}

// Close closes the WebSocket connection and stops the ping handler
func (c *Client) Close() error {
	select {
	case <-c.done:
		return nil
	default:
		close(c.done)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = false
	if c.ws != nil {
		err := c.ws.Close()
		c.ws = nil
		return err
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
