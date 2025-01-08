# TradingView WebSocket Client

A Go library for receiving real-time market data from TradingView's WebSocket API.

## Features

- Real-time market data streaming with comprehensive quote data
- Support for multiple symbols and exchanges (NASDAQ, BINANCE, HKEX, etc.)
- Automatic session management and WebSocket connection handling
- Rich market data fields including:
  - Last price and timestamp
  - Price change and percentage
  - Trading volume
  - Currency information
  - Exchange details
  - Session status
- Dynamic symbol management (add/remove symbols during runtime)
- Structured logging with slog
- Configurable client options

## Installation

```bash
go get github.com/iiiyu/tradingview-ws-client
```

## Usage

### Basic Example

```go
package main

import (
    "github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

func main() {
    // Create a new client with your TradingView auth token
    client, err := tvwsclient.NewClient("your-auth-token")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Create data channel for receiving updates
    dataChan := make(chan map[string]interface{})

    // Start receiving data
    go client.ReadMessage(dataChan)

    // Setup quote session
    qsSession := tvwsclient.GenerateSession("qs_")
    tvwsclient.SendQuoteCreateSessionMessage(client, qsSession)
    tvwsclient.SendQuoteSetFieldsMessage(client, qsSession)

    // Add symbols to watch
    symbols := []string{"NASDAQ:AAPL", "BINANCE:BTCUSDT"}
    tvwsclient.SendQuoteAddSymbolsMessage(client, qsSession, symbols)

    // Process incoming data
    for data := range dataChan {
        if response, ok := data["p"].([]interface{}); ok && len(response) >= 2 {
            if quote, ok := response[1].(tvwsclient.QuoteData); ok {
                // Handle quote data
                _ = quote.Name        // Symbol name
                _ = quote.Values.LastPrice   // Current price
                _ = quote.Values.Change     // Price change
                _ = quote.Values.Volume     // Trading volume
            }
        }
    }
}
```

### Configuration

The client supports configuration through YAML files:

```yaml
auth:
  token: "your-tradingview-auth-token"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
