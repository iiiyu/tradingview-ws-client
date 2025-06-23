package tvwsclient

import (
	"context"
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

// Add a new ConnectionState type for better state management
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
)

// Client represents a TradingView WebSocket client
type Client struct {
	ws            *websocket.Conn
	mu            sync.Mutex // protects ws
	requestHeader http.Header
	wsURL         string
	maxRetries    int
	done          chan struct{} // Channel to signal connection close
	pingInterval  time.Duration // Interval for sending ping messages
	writeTimeout  time.Duration // Timeout for write operations
	readTimeout   time.Duration // Timeout for read operations
	state         ConnectionState
	cancel        context.CancelFunc
	
	// Reconnection state
	reconnecting  bool
	retryCount    int
	lastConnectTime time.Time
	
	// Callback for handling reconnection events
	onReconnect   func() error
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
		state:         StateDisconnected, // Initial state
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
	c.state = StateConnecting // Set state to connecting

	if c.ws != nil {
		c.ws.Close()
		c.ws = nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(c.wsURL, c.requestHeader)
	if err != nil {
		c.state = StateDisconnected // Set state back to disconnected on failure
		c.mu.Unlock()
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.ws = conn
	c.state = StateConnected  // Set state to connected after successful connection

	// Setup ping handler to respond to server pings
	c.ws.SetPingHandler(func(appData string) error {
		c.mu.Lock()
		defer c.mu.Unlock()
		if c.ws == nil {
			return fmt.Errorf("connection is closed")
		}
		return c.ws.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(10*time.Second))
	})
	
	// Release mutex before calling SendInitMessage to avoid deadlock
	c.mu.Unlock()

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
			// Check connection state before sending ping
			if c.ws == nil || c.state != StateConnected {
				c.mu.Unlock()
				slog.Warn("skipping ping - connection not ready", "state", c.state)
				continue
			}
			
			err := c.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(c.writeTimeout))
			c.mu.Unlock()

			if err != nil {
				slog.Error("ping error", "error", err, "retry_count", c.retryCount)
				
				// Don't attempt reconnection if already reconnecting
				c.mu.Lock()
				shouldReconnect := !c.reconnecting && c.state != StateConnecting
				c.mu.Unlock()
				
				if shouldReconnect {
					go func() {
						if err := c.reconnect(); err != nil {
							slog.Error("failed to reconnect after ping error", "error", err)
						}
					}()
				}
			} else {
				// Reset retry count on successful ping
				c.mu.Lock()
				if c.retryCount > 0 {
					slog.Debug("ping successful, resetting retry count", "previous_retries", c.retryCount)
					c.retryCount = 0
				}
				c.mu.Unlock()
			}
		case <-c.done:
			return
		}
	}
}

// reconnect attempts to reconnect to the WebSocket server
func (c *Client) reconnect() error {
	c.mu.Lock()
	
	// Prevent multiple reconnection attempts
	if c.reconnecting {
		c.mu.Unlock()
		return fmt.Errorf("reconnection already in progress")
	}
	
	c.reconnecting = true
	c.state = StateConnecting
	
	// Close existing connection if any
	if c.ws != nil {
		c.ws.Close()
		c.ws = nil
	}
	
	c.mu.Unlock()
	
	defer func() {
		c.mu.Lock()
		c.reconnecting = false
		c.mu.Unlock()
	}()
	
	// Exponential backoff with jitter
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		c.retryCount++
		
		// Calculate delay: 2^attempt seconds, capped at 30 seconds
		delay := time.Duration(1<<uint(attempt)) * time.Second
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}
		
		// Add jitter (±25% of delay)
		jitter := time.Duration(float64(delay) * 0.25 * (2*float64(time.Now().UnixNano()%2) - 1))
		delay += jitter
		
		if attempt > 0 {
			slog.Info("reconnection attempt", 
				"attempt", attempt+1, 
				"max_retries", c.maxRetries,
				"delay", delay)
			time.Sleep(delay)
		}
		
		// Attempt to connect
		err := c.connect()
		if err != nil {
			slog.Error("reconnection failed", 
				"attempt", attempt+1, 
				"error", err)
			continue
		}
		
		// Connection successful
		c.mu.Lock()
		c.state = StateConnected
		c.retryCount = 0
		c.lastConnectTime = time.Now()
		c.mu.Unlock()
		
		slog.Info("reconnection successful", "attempt", attempt+1)
		
		// Call reconnection callback if set
		if c.onReconnect != nil {
			if err := c.onReconnect(); err != nil {
				slog.Error("reconnection callback failed", "error", err)
			}
		}
		
		return nil
	}
	
	// All attempts failed
	c.mu.Lock()
	c.state = StateDisconnected
	c.mu.Unlock()
	
	return fmt.Errorf("failed to reconnect after %d attempts", c.maxRetries)
}

func (c *Client) Reconnect() error {
	return c.reconnect()
}

// SetReconnectCallback sets a callback function to be called after successful reconnection
func (c *Client) SetReconnectCallback(callback func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onReconnect = callback
}

// GetConnectionState returns the current connection state
func (c *Client) GetConnectionState() ConnectionState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// IsConnected returns true if the client is currently connected
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state == StateConnected && c.ws != nil
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
		if c.ws == nil || c.state != StateConnected {
			shouldAttemptReconnect := !c.reconnecting && retries < c.maxRetries
			c.mu.Unlock()
			
			if shouldAttemptReconnect {
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
				if retries >= c.maxRetries {
					return WrapConnectionError("read_message.max_retries", ErrReconnectFailed)
				}
				// If reconnection is already in progress, wait a bit
				time.Sleep(100 * time.Millisecond)
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

			// Check if we should attempt reconnect
			c.mu.Lock()
			shouldAttemptReconnect := !c.reconnecting && retries < c.maxRetries
			c.mu.Unlock()
			
			if shouldAttemptReconnect {
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
			} else {
				if retries >= c.maxRetries {
					return fmt.Errorf("error reading message after %d reconnection attempts: %w", retries, err)
				}
				// If reconnection is already in progress, wait a bit
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		// Reset retry counter on successful message
		retries = 0

		// Handle heartbeat messages
		if heartbeatRegex.Match(message) {
			c.mu.Lock()
			if c.ws == nil {
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
	c.cancel() // Cancel context

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateDisconnected {
		return nil
	}

	c.state = StateDisconnected
	if c.ws != nil {
		if err := c.ws.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(time.Second),
		); err != nil {
			slog.Error("error sending close message", "error", err)
		}

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
