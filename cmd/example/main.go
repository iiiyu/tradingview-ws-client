package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Create a new client
	client, err := tvwsclient.NewClient()
	if err != nil {
		slog.Error("failed to create client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	// Example symbols
	symbols := []string{
		"NASDAQ:AAPL",
		"NASDAQ:MSFT",
		"NASDAQ:GOOGL",
		"NASDAQ:AMZN",
		"NASDAQ:META",
		"BINANCE:BTCUSDT",
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create data channel
	dataChan := make(chan map[string]interface{})

	// Start receiving data
	go func() {
		if err := client.GetLatestTradeInfo(symbols, dataChan); err != nil {
			slog.Error("failed to get trade info", "error", err)
		}
	}()

	slog.Info("starting to receive trade data", "symbols", symbols)

	// Main loop
	for {
		select {
		case data := <-dataChan:
			slog.Debug("received trade data", "data", data)
		case <-sigChan:
			slog.Info("shutting down...")
			return
		}
	}
}
