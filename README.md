# TradingView WebSocket Client

A Go library for receiving real-time market data from TradingView's WebSocket API.

## Features

- Real-time market data streaming with comprehensive quote data
- Support for multiple symbols and exchanges (NASDAQ, BINANCE, HKEX, etc.)
- Automatic WebSocket connection handling with reconnection support
- Configurable client options including:
  - Custom WebSocket URL
  - Custom headers
  - Maximum retry attempts
- Rich market data fields including:
  - Last price and timestamp
  - Price change and percentage
  - Trading volume
  - Currency information
  - Exchange details
  - Session status
- Dynamic symbol management (add/remove symbols during runtime)
- Structured logging with slog
- Database integration support (PostgreSQL)

## Installation

```bash
go get github.com/iiiyu/tradingview-ws-client
```

## Configuration

The client supports configuration through environment variables and flags:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=tradingview
DB_SSLMODE=disable

# TradingView Configuration
TRADINGVIEW_AUTH_TOKEN=your-auth-token
```

## Usage

### Basic Example

```go
package main

import (
    "log"
    tvws "github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

func main() {
    // Create a new client with options
    client, err := tvws.NewClient(
        "your-auth-token",
        // Optional: Configure custom options
        // tvws.WithMaxRetries(5),
        // tvws.WithCustomHeaders(headers),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create data channel for receiving updates
    dataChan := make(chan tvws.TVResponse)

    // Start receiving data in a goroutine
    go func() {
        if err := client.ReadMessage(dataChan); err != nil {
            log.Printf("Error reading messages: %v", err)
        }
    }()

    // Setup quote session
    qsSession := tvws.GenerateSession("qs_")
    if err := tvws.SendQuoteCreateSessionMessage(client, qsSession); err != nil {
        log.Fatal(err)
    }
    if err := tvws.SendQuoteSetFieldsMessage(client, qsSession); err != nil {
        log.Fatal(err)
    }

    // Add symbols to watch
    symbols := []string{"NASDAQ:AAPL", "BINANCE:BTCUSDT"}
    if err := tvws.SendQuoteAddSymbolsMessage(client, qsSession, symbols); err != nil {
        log.Fatal(err)
    }

    // Process incoming data
    for data := range dataChan {
        // Handle the data based on your needs
        log.Printf("Received data: %+v", data)
    }
}
```

### Error Handling

The client includes built-in error handling and reconnection logic. You can customize the retry behavior using the `WithMaxRetries` option when creating the client.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
