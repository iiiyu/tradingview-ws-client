# TradingView WebSocket Client

A robust Go library for connecting to TradingView's WebSocket API to receive real-time market data.

## 🚀 Features

### Core WebSocket Client
- **Real-time market data streaming** with comprehensive quote and candle data
- **Robust reconnection strategy** with exponential backoff and jitter
- **Thread-safe operations** with comprehensive mutex protection
- **Automatic session restoration** after connection interruptions
- **Connection state management** with proper state transitions
- **Graceful error recovery** with structured retry logic
- **Support for multiple data types** (quotes, candles, level2, trades)

### Reliability & Stability
- **Deadlock prevention** with careful lock management
- **Concurrent access protection** for all WebSocket send operations
- **Enhanced ping/pong handling** for connection health monitoring
- **Configurable timeouts and retry limits**
- **Structured error handling** with context preservation

## 📦 Installation

```bash
go get github.com/iiiyu/tradingview-ws-client
```

## ⚡ Quick Start

### Basic Usage

```go
package main

import (
    "log"
    "log/slog"
    "time"
    
    tvws "github.com/iiiyu/tradingview-ws-client/tvwsclient"
)

func main() {
    // Initialize AuthTokenManager (required)
    httpClient := tvws.NewTVHttpClient(
        "https://www.tradingview.com",
        "your_device_token",
        "your_session_id", 
        "your_session_sign",
    )
    tvws.InitAuthTokenManager(httpClient)

    // Create WebSocket client
    client, err := tvws.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create data channel
    dataChan := make(chan tvws.TVResponse, 1000)

    // Start reading messages
    go func() {
        if err := client.ReadMessage(dataChan); err != nil {
            slog.Error("Error reading messages", "error", err)
        }
    }()

    // Subscribe to quotes
    session := tvws.GenerateSession("qs_")
    if err := tvws.SubscriptionQuoteSessionSymbol(client, session, "NASDAQ:AAPL"); err != nil {
        log.Fatal(err)
    }

    // Process incoming data
    for data := range dataChan {
        slog.Info("Received data", "method", data.Method, "data", data)
    }
}
```

### Advanced Configuration

```go
// Create client with custom options
client, err := tvws.NewClient(
    tvws.WithMaxRetries(10),
    tvws.WithPingInterval(30*time.Second),
    tvws.WithWriteTimeout(60*time.Second),
    tvws.WithReadTimeout(60*time.Second),
)
if err != nil {
    log.Fatal(err)
}

// Set up reconnection callback for session restoration
client.SetReconnectCallback(func() error {
    slog.Info("Connection restored, resubscribing...")
    // Your session restoration logic here
    return nil
})

// Check connection state
if client.IsConnected() {
    slog.Info("Client is connected")
}

// Get current connection state
state := client.GetConnectionState()
switch state {
case tvws.StateConnected:
    slog.Info("Connection is healthy")
case tvws.StateConnecting:
    slog.Info("Connection in progress")
case tvws.StateDisconnected:
    slog.Warn("Connection lost")
}
```

## 📡 API Reference

### Client Creation

```go
// Create client with default settings
client, err := tvws.NewClient()

// Create client with options
client, err := tvws.NewClient(
    tvws.WithMaxRetries(5),
    tvws.WithPingInterval(30*time.Second),
)
```

### Connection Management

```go
// Check if connected
isConnected := client.IsConnected()

// Get connection state
state := client.GetConnectionState()

// Manual reconnection
err := client.Reconnect()

// Set reconnection callback
client.SetReconnectCallback(func() error {
    // Handle reconnection
    return nil
})

// Close connection
client.Close()
```

### Data Subscriptions

```go
// Subscribe to quote data
session := tvws.GenerateSession("qs_")
err := tvws.SubscriptionQuoteSessionSymbol(client, session, "NASDAQ:AAPL")

// Subscribe to candle data
session := tvws.GenerateSession("cs_")
err := tvws.SubscriptionChartSessionSymbol(client, session, "NASDAQ:AAPL", "1D", 300)

// Unsubscribe from quotes
err := tvws.SendQuoteRemoveSymbolsMessage(client, session, []string{"NASDAQ:AAPL"})

// Unsubscribe from candles
err := tvws.SendChartDeleteSessionMessage(client, session)
```

### Message Processing

```go
// Create message router for handling different message types
router := tvws.NewMessageRouter(logger)

// Register handlers
router.RegisterHandler(tvws.MethodQuoteData, quoteHandler)
router.RegisterHandler(tvws.MethodTimescaleUpdate, candleHandler)

// Route messages
for data := range dataChan {
    if err := router.RouteMessage(context.Background(), data); err != nil {
        slog.Error("Failed to route message", "error", err)
    }
}
```

## 🔧 Configuration

### Authentication Setup

Before using the client, you must initialize the AuthTokenManager:

```go
httpClient := tvws.NewTVHttpClient(
    "https://www.tradingview.com",
    "your_device_token",
    "your_session_id",
    "your_session_sign",
)
tvws.InitAuthTokenManager(httpClient)
```

### Client Options

```go
type Option func(*Client)

// Set maximum retry attempts for reconnection
func WithMaxRetries(retries int) Option

// Set ping interval for connection health checks
func WithPingInterval(interval time.Duration) Option

// Set write timeout for WebSocket operations
func WithWriteTimeout(timeout time.Duration) Option

// Set read timeout for WebSocket operations  
func WithReadTimeout(timeout time.Duration) Option
```

## 🔄 Connection States

The client manages three connection states:

- **StateDisconnected**: No active connection
- **StateConnecting**: Connection attempt in progress
- **StateConnected**: Active and healthy connection

## 📊 Data Types

### Quote Data
```go
type QuoteDataMessage struct {
    Data struct {
        Name   string `json:"n"`
        Values struct {
            Change        float64 `json:"ch"`
            LastPrice     float64 `json:"lp"`
            LastPriceTime int64   `json:"lp_time"`
            Volume        float64 `json:"volume"`
            Bid           float64 `json:"bid"`
            Ask           float64 `json:"ask"`
            // ... more fields
        } `json:"v"`
    } `json:"p"`
}
```

### Candle Data
```go
type TimescaleUpdateMessage struct {
    ChartSessionID string `json:"p"`
    Data struct {
        SDS1 struct {
            S []struct {
                V []float64 `json:"v"` // [timestamp, open, high, low, close, volume]
            } `json:"s"`
        } `json:"sds_1"`
    } `json:"p"`
}
```

## 🚨 Error Handling

The library provides structured error types:

```go
// Check for specific error types
if tvws.IsConnectionError(err) {
    // Handle connection errors
}

if tvws.IsRetryableError(err) {
    // Handle retryable errors
}

// Wrap errors with context
err = tvws.WrapConnectionError("operation_name", originalErr)
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## 📋 Requirements

- **Go 1.24+**
- **TradingView credentials** (device token, session ID, session signature)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run tests and ensure they pass
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [TradingView HTTP Service](https://github.com/iiiyu/tradingview-http-service) - HTTP API service using this WebSocket client

## 🆘 Support

For issues and questions:

1. Check the [examples](./examples/) directory
2. Review the API documentation above
3. Open an issue with detailed information