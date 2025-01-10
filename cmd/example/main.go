package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	// symbols := []string{
	// 	"NASDAQ:AAPL",
	// 	// "BINANCE:BTCUSDT",
	// 	// "HKEX:700",
	// 	// "HKEX_DLY:1810",
	// }
	// slog.Debug("starting to receive trade data", "symbols", symbols)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create data channel
	dataChan := make(chan tvwsclient.TVResponse)

	// Start receiving data
	go func() {
		if err := client.ReadMessage(dataChan); err != nil {
			slog.Error("failed to get trade info", "error", err)
		}
	}()

	go func() {
		qsSession := tvwsclient.GenerateSession("qs_")
		csSession := tvwsclient.GenerateSession("cs_")
		csSymbol := "BINANCE:BTCUSDT"

		if err := tvwsclient.SendQuoteCreateSessionMessage(client, qsSession); err != nil {
			slog.Error("failed to send quote create session message ", "error", err)
		}

		if err := tvwsclient.SendQuoteSetFieldsMessage(client, qsSession); err != nil {
			slog.Error("failed to send quote set fields session message ", "error", err)
		}

		if err := tvwsclient.SubscriptionChartSessionSymbol(client, csSession, csSymbol, "10S", 10); err != nil {
			slog.Error("failed to subscription chart session", "error", err)
		}

		csSession2 := tvwsclient.GenerateSession("cs_")
		csSymbol2 := "BINANCE:BTCUSDT"
		if err := tvwsclient.SubscriptionChartSessionSymbol(client, csSession2, csSymbol2, "30S", 20); err != nil {
			slog.Error("failed to subscription chart session", "error", err)
		}

		// if err := tvwsclient.SendChartDeleteSessionMessage(client, csSession); err != nil {
		// 	slog.Error("failed to send chart delete session message", "error", err)
		// }

		// if err := tvwsclient.SendChartCreateSessionMessage(client, csSession); err != nil {
		// 	slog.Error("failed to send chart create session message ", "error", err)
		// }

		// if err := tvwsclient.SendSwitchTimezoneMessage(client, csSession); err != nil {
		// 	slog.Error("failed to send switch timezone message ", "error", err)
		// }

		// if err := tvwsclient.SendResolveSymbolMessage(client, csSession, csSymbol); err != nil {
		// 	slog.Error("failed to send resolve symbol message ", "error", err)
		// }

		// if err := tvwsclient.SendCreateSeriesMessage(client, csSession, "1"); err != nil {
		// 	slog.Error("failed to send create series message ", "error", err)
		// }
		// if err := tvwsclient.SendQuoteAddSymbolsMessage(client, qsSession, symbols); err != nil {
		// 	slog.Error("failed to send add quote symbols session message ", "error", err)
		// }
		// if err := tvwsclient.SendQuoteAddSymbolsMessage(client, qsSession, []string{"BINANCE:SOLUSDT", "BINANCE:ETHUSDT"}); err != nil {
		// 	slog.Error("failed to send add quote symbols session message ", "error", err)
		// }

		// time.Sleep(10 * time.Second)
		// if err := tvwsclient.SendQuoteRemoveSymbolsMessage(client, qsSession, []string{"BINANCE:BTCUSDT"}); err != nil {
		// 	slog.Error("failed to send remove quote symbols session message ", "error", err)
		// }
	}()

	// Main loop
	for {
		select {
		case data := <-dataChan:
			slog.Debug("data", "data.Method", data.Method)
			slog.Debug("data", "data", data.Params)
			switch data.Method {
			case "qsd":
				// please make the data.Param to tvwsclient.QuoteData
				if len(data.Params) >= 2 {
					// Convert the interface{} back to JSON
					paramJSON, err := json.Marshal(data.Params[1])
					if err != nil {
						// return nil, fmt.Errorf("failed to marshal param: %w", err)
						slog.Error("failed to marshal param", "error", err)
					}

					// Unmarshal into QuoteData
					var quoteData tvwsclient.QuoteData
					if err := json.Unmarshal(paramJSON, &quoteData); err != nil {
						slog.Error("failed to unmarshal quote data", "error", err)
					}
					// if quote, ok := data.Params[1].(tvwsclient.QuoteData); ok {
					slog.Info("received quote data",
						"symbol", quoteData.Name,
						"price", quoteData.Values.LastPrice,
						"change", quoteData.Values.Change,
						"volume", quoteData.Values.Volume,
					)
					// }
				}
			}
		// if response, ok := data["p"].([]interface{}); ok && len(response) >= 2 {
		// 	switch data["m"] {
		// 	case "qsd":
		// 		if quote, ok := response[1].(tvwsclient.QuoteData); ok {
		// 			slog.Info("received quote data",
		// 				"symbol", quote.Name,
		// 				"price", quote.Values.LastPrice,
		// 				"change", quote.Values.Change,
		// 				"volume", quote.Values.Volume,
		// 			)
		// 		}
		// 	case "series_loading":
		// 		if msg, ok := response[1].(tvwsclient.SeriesLoadingMessage); ok {
		// 			slog.Info("received series loading",
		// 				"session", response[0],
		// 				"series", msg,
		// 			)
		// 		}
		// 	case "symbol_resolved":
		// 		if msg, ok := response[1].(tvwsclient.SymbolResolvedMessage); ok {
		// 			slog.Info("received symbol resolved",
		// 				"session", response[0],
		// 				"message", msg,
		// 			)
		// 		}
		// 	case "timescale_update":
		// 		if msg, ok := response[1].(tvwsclient.TimescaleUpdateMessage); ok {
		// 			slog.Info("received timescale update",
		// 				"session", response[0],
		// 				"message", msg,
		// 			)
		// 		}
		// 	case "series_completed":
		// 		if msg, ok := response[1].(tvwsclient.SeriesCompletedMessage); ok {
		// 			slog.Info("received series completed",
		// 				"session", response[0],
		// 				"message", msg,
		// 			)
		// 		}
		// 	}
		// }

		case <-sigChan:
			slog.Info("shutting down...")
			return
		}
	}
}
