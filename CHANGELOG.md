# Changelog

## [0.1.0] - 2025-06-23

### Initial Release

This is the first release of the TradingView WebSocket Client library, extracted from the original monolithic application.

### Features

#### Core WebSocket Client
- **Real-time WebSocket connection** to TradingView's data streams
- **Robust reconnection strategy** with exponential backoff and jitter
- **Thread-safe operations** with comprehensive mutex protection
- **Connection state management** (Disconnected, Connecting, Connected)
- **Configurable timeouts and retry limits**
- **Support for multiple data types** (quotes, candles, level2, trades)

#### Reliability & Stability
- **Automatic session restoration** after connection interruptions
- **Deadlock prevention** with careful lock management
- **Concurrent access protection** for all WebSocket send operations
- **Enhanced ping/pong handling** for connection health monitoring
- **Graceful error recovery** with structured retry logic
- **Nil pointer protection** for all WebSocket operations

#### Message Processing
- **Message router pattern** for extensible message handling
- **Structured error handling** with context preservation
- **Support for quote data, candle updates, and data updates**
- **Configurable message buffering** with channels

#### Authentication
- **TradingView session management** with device tokens
- **Automatic token refresh** and session validation
- **HTTP client integration** for authentication flows

### API

#### Client Creation
```go
client, err := tvwsclient.NewClient(
    tvwsclient.WithMaxRetries(5),
    tvwsclient.WithPingInterval(30*time.Second),
)
```

#### Connection Management
```go
client.IsConnected()
client.GetConnectionState()
client.Reconnect()
client.SetReconnectCallback(callback)
client.Close()
```

#### Data Subscriptions
```go
tvwsclient.SubscriptionQuoteSessionSymbol(client, session, "NASDAQ:AAPL")
tvwsclient.SubscriptionChartSessionSymbol(client, session, "NASDAQ:AAPL", "1D", 300)
```

### Breaking Changes

This is a new library extracted from a larger application, so there are no breaking changes from previous versions.

### Migration from Original Codebase

If migrating from the original `github.com/iiiyu/tradingview-ws-client` monolithic repository:

1. Update import paths:
   ```go
   // Old
   "github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
   
   // New  
   "github.com/iiiyu/tradingview-ws-client/tvwsclient"
   ```

2. Initialize AuthTokenManager before creating client:
   ```go
   httpClient := tvwsclient.NewTVHttpClient(baseURL, deviceToken, sessionID, sessionSign)
   tvwsclient.InitAuthTokenManager(httpClient)
   ```

3. Use new connection state methods:
   ```go
   if client.IsConnected() {
       // Connection is healthy
   }
   ```

### Dependencies

- `github.com/gorilla/websocket v1.5.0`

### Documentation

- [README.md](./README.md) - Complete usage guide
- [examples/basic/](./examples/basic/) - Basic usage example

### Contributors

- Initial extraction and refactoring of WebSocket client functionality
- Implementation of robust reconnection strategy
- Thread-safety improvements and deadlock prevention
- Comprehensive error handling and state management