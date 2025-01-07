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

// TVResponse represents the top-level response structure
type TVResponse struct {
	Method string        `json:"m"` // "qsd" for quote data
	Params []interface{} `json:"p"` // Array containing session ID and quote data
}

// QuoteData represents the structure of quote data for a symbol
type QuoteData struct {
	Name   string     `json:"n"` // Symbol name (e.g., "BINANCE:BTCUSDT")
	Status string     `json:"s"` // Status ("ok")
	Values SymbolData `json:"v"` // Actual symbol data
}

// SymbolData represents the trading data for a symbol
type SymbolData struct {
	Change         float64 `json:"ch"`              // Price change
	ChangePercent  float64 `json:"chp"`             // Price change percentage
	CurrencyCode   string  `json:"currency_code"`   // Currency (USD, USDT, etc.)
	CurrentSession string  `json:"current_session"` // Session status
	Description    string  `json:"description"`     // Full name
	Exchange       string  `json:"exchange"`        // Exchange name
	IsTradable     bool    `json:"is_tradable"`
	LastPrice      float64 `json:"lp"`         // Last price
	LastPriceTime  int64   `json:"lp_time"`    // Timestamp
	PriceScale     int     `json:"pricescale"` // Price scale factor
	ShortName      string  `json:"short_name"` // Symbol short name
	Type           string  `json:"type"`       // Asset type (stock, spot, etc.)
	Volume         float64 `json:"volume"`     // Trading volume
}
