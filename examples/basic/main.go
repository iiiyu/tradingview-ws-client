package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tvws "github.com/iiiyu/tradingview-ws-client/tvwsclient"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize AuthTokenManager (required)
	// You need to get these from your TradingView account
	deviceToken := os.Getenv("TRADINGVIEW_DEVICE_TOKEN")
	sessionID := os.Getenv("TRADINGVIEW_SESSION_ID")
	sessionSign := os.Getenv("TRADINGVIEW_SESSION_SIGN")

	if deviceToken == "" || sessionID == "" || sessionSign == "" {
		log.Fatal("Please set TRADINGVIEW_DEVICE_TOKEN, TRADINGVIEW_SESSION_ID, and TRADINGVIEW_SESSION_SIGN environment variables")
	}

	httpClient := tvws.NewTVHttpClient(
		"https://www.tradingview.com",
		deviceToken,
		sessionID,
		sessionSign,
	)
	tvws.InitAuthTokenManager(httpClient)

	// Create WebSocket client
	client, err := tvws.NewClient()
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}
	defer client.Close()

	// Set up reconnection callback
	client.SetReconnectCallback(func() error {
		slog.Info("Connection restored successfully")
		return nil
	})

	// Create data channel
	dataChan := make(chan tvws.TVResponse, 1000)

	// Start reading messages in a goroutine
	go func() {
		defer close(dataChan)
		if err := client.ReadMessage(dataChan); err != nil {
			slog.Error("Error reading messages", "error", err)
		}
	}()

	// Subscribe to Apple stock quotes
	quoteSession := tvws.GenerateSession("qs_")
	if err := tvws.SubscriptionQuoteSessionSymbol(client, quoteSession, "NASDAQ:AAPL"); err != nil {
		log.Fatal("Failed to subscribe to AAPL quotes:", err)
	}
	slog.Info("Subscribed to NASDAQ:AAPL quotes", "session", quoteSession)

	// Subscribe to Bitcoin candle data
	candleSession := tvws.GenerateSession("cs_")
	if err := tvws.SubscriptionChartSessionSymbol(client, candleSession, "BINANCE:BTCUSDT", "1D", 100); err != nil {
		log.Fatal("Failed to subscribe to BTCUSDT candles:", err)
	}
	slog.Info("Subscribed to BINANCE:BTCUSDT 1D candles", "session", candleSession)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Process incoming data
	slog.Info("Starting to process market data... Press Ctrl+C to stop")
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				slog.Info("Data channel closed")
				return
			}

			// Log received data based on method
			switch data.Method {
			case tvws.MethodQuoteData:
				slog.Info("Received quote data", 
					"method", data.Method,
					"data_type", "quote")
			case tvws.MethodTimescaleUpdate:
				slog.Info("Received candle update", 
					"method", data.Method,
					"data_type", "candle")
			case tvws.MethodDataUpdate:
				slog.Info("Received data update", 
					"method", data.Method,
					"data_type", "update")
			default:
				slog.Debug("Received other message", 
					"method", data.Method)
			}

		case <-sigChan:
			slog.Info("Received shutdown signal, cleaning up...")
			
			// Clean shutdown
			if client.IsConnected() {
				// Unsubscribe from sessions
				if err := tvws.SendQuoteRemoveSymbolsMessage(client, quoteSession, []string{"NASDAQ:AAPL"}); err != nil {
					slog.Error("Failed to unsubscribe from quotes", "error", err)
				}
				if err := tvws.SendChartDeleteSessionMessage(client, candleSession); err != nil {
					slog.Error("Failed to unsubscribe from candles", "error", err)
				}
			}
			
			return
		}
	}
}