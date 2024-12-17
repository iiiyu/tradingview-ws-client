package tvwsclient

import (
	"net/http"
)

// Option represents a client option
type Option func(*Client)

// WithCustomHeaders sets custom HTTP headers for the WebSocket connection
func WithCustomHeaders(headers http.Header) Option {
	return func(c *Client) {
		c.requestHeader = headers
	}
}

// WithCustomURL sets a custom WebSocket URL
func WithCustomURL(url string) Option {
	return func(c *Client) {
		c.wsURL = url
	}
}

// TradeData represents the real-time trade information
type TradeData struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"lp"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"lp_time"`
	// Add more fields as needed
}
