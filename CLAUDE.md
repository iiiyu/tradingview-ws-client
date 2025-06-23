# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Run specific test files
go test ./tvwsclient/auth_token_test.go
go test ./tvwsclient/utils_test.go
```

### Building and Running
```bash
# Build the module
go build ./...

# Run the basic example
cd examples/basic
go run main.go

# Build and run with environment variables
TRADINGVIEW_DEVICE_TOKEN=your_token TRADINGVIEW_SESSION_ID=your_id TRADINGVIEW_SESSION_SIGN=your_sign go run main.go
```

### Module Management
```bash
# Update dependencies
go mod tidy

# Add new dependencies
go get package_name

# Verify module
go mod verify
```

## Architecture Overview

### Core Components

**Client (`tvwsclient/client.go`)**: The main WebSocket client that manages connections to TradingView's WebSocket API. Features robust reconnection strategy with exponential backoff, thread-safe operations, and connection state management (StateDisconnected, StateConnecting, StateConnected).

**Authentication (`tvwsclient/auth_token.go`)**: Manages TradingView authentication tokens with automatic refresh capabilities. Must be initialized before creating clients using `InitAuthTokenManager()`.

**Message Router (`tvwsclient/message_router.go`)**: Routes incoming WebSocket messages to appropriate handlers based on message method types (qsd, timescale_update, du, etc.).

**Interfaces (`tvwsclient/interfaces.go`)**: Defines comprehensive interfaces for client operations, message handling, authentication, and data persistence. Supports dependency injection and testing.

### Message Flow

1. **Authentication**: Initialize `AuthTokenManager` with TradingView credentials
2. **Connection**: Create client with `NewClient()` and optional configuration
3. **Subscription**: Subscribe to quote/chart data using session-based APIs
4. **Processing**: Route incoming messages through `MessageRouter` to registered handlers
5. **Data Handling**: Process quote data, candle data, and other market information

### Key Patterns

**Session Management**: Uses session-based subscriptions with generated IDs (e.g., `qs_` for quotes, `cs_` for charts).

**Error Handling**: Structured error types with `IsConnectionError()` and `IsRetryableError()` helpers for specific error handling.

**Concurrency**: Thread-safe operations with mutex protection, especially for WebSocket send operations.

**Reconnection**: Automatic reconnection with callbacks for session restoration after connection interruptions.

## Required Environment Variables

For examples and testing:
- `TRADINGVIEW_DEVICE_TOKEN`: Your TradingView device token
- `TRADINGVIEW_SESSION_ID`: Your TradingView session ID  
- `TRADINGVIEW_SESSION_SIGN`: Your TradingView session signature

## Common Development Tasks

### Adding New Message Types
1. Add method constant to `types.go`
2. Define message struct in `types.go`
3. Add handler method to `MessageHandler` interface
4. Implement routing logic in `message_router.go`

### Testing WebSocket Functionality
Use the basic example in `examples/basic/` as a starting point. Set required environment variables and modify the symbol subscriptions as needed.

### Extending Client Options
Add new option functions following the `WithMaxRetries()` pattern in the client constructor.