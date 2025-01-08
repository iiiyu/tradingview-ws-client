package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// Setup command line flags
	configPath := pflag.StringP("config", "c", "config.yaml", "path to config file")
	pflag.Parse()

	// Setup viper
	viper.SetConfigFile(*configPath)
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", "error", err)
		os.Exit(1)
	}

	// Get auth token from config
	authToken := viper.GetString("auth.token")
	if authToken == "" {
		slog.Error("auth token not found in config")
		os.Exit(1)
	}

	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	// Create a new client
	client, err := tvwsclient.NewClient(authToken)
	if err != nil {
		slog.Error("failed to create client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	// Example symbols
	symbols := []string{
		"NASDAQ:AAPL",
		// "BINANCE:BTCUSDT",
		// "HKEX:700",
		// "HKEX_DLY:1810",
	}
	slog.Info("starting to receive trade data", "symbols", symbols)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create data channel
	dataChan := make(chan map[string]interface{})

	// Start receiving data
	go func() {
		if err := client.ReadMessage(dataChan); err != nil {
			slog.Error("failed to get trade info", "error", err)
		}
	}()

	go func() {
		qsSession := tvwsclient.GenerateSession("qs_")
		if err := tvwsclient.SendQuoteCreateSessionMessage(client, qsSession); err != nil {
			slog.Error("failed to send quote create session message ", "error", err)
		}
		if err := tvwsclient.SendQuoteSetFieldsMessage(client, qsSession); err != nil {
			slog.Error("failed to send quote set fields session message ", "error", err)
		}
		if err := tvwsclient.SendQuoteAddSymbolsMessage(client, qsSession, symbols); err != nil {
			slog.Error("failed to send add quote symbols session message ", "error", err)
		}
		// if err := tvwsclient.SendQuoteAddSymbolsMessage(client, qsSession, []string{"BINANCE:SOLUSDT", "BINANCE:ETHUSDT"}); err != nil {
		// 	slog.Error("failed to send add quote symbols session message ", "error", err)
		// }

		time.Sleep(10 * time.Second)
		if err := tvwsclient.SendQuoteRemoveSymbolsMessage(client, qsSession, []string{"BINANCE:BTCUSDT"}); err != nil {
			slog.Error("failed to send remove quote symbols session message ", "error", err)
		}
	}()

	// Main loop
	for {
		select {
		case data := <-dataChan:
			if response, ok := data["p"].([]interface{}); ok && len(response) >= 2 {
				if quote, ok := response[1].(tvwsclient.QuoteData); ok {
					slog.Info("received trade data",
						"symbol", quote.Name,
						"price", quote.Values.LastPrice,
						"change", quote.Values.Change,
						"volume", quote.Values.Volume,
					)
				}
			}
		case <-sigChan:
			slog.Info("shutting down...")
			return
		}
	}
}
