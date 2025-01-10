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
	symbols := []string{
		"NASDAQ:AAPL",
		"BINANCE:BTCUSDT",
		// "HKEX:700",
		// "HKEX_DLY:1810",
	}
	slog.Debug("starting to receive trade data", "symbols", symbols)

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
		// csSession := tvwsclient.GenerateSession("cs_")
		// csSymbol := "BINANCE:BTCUSDT"

		if err := tvwsclient.SendQuoteCreateSessionMessage(client, qsSession); err != nil {
			slog.Error("failed to send quote create session message ", "error", err)
		}

		if err := tvwsclient.SendQuoteSetFieldsMessage(client, qsSession); err != nil {
			slog.Error("failed to send quote set fields session message ", "error", err)
		}

		// if err := tvwsclient.SubscriptionChartSessionSymbol(client, csSession, csSymbol, "10S", 10); err != nil {
		// 	slog.Error("failed to subscription chart session", "error", err)
		// }

		// csSession2 := tvwsclient.GenerateSession("cs_")
		// csSymbol2 := "BINANCE:BTCUSDT"
		// if err := tvwsclient.SubscriptionChartSessionSymbol(client, csSession2, csSymbol2, "30S", 20); err != nil {
		// 	slog.Error("failed to subscription chart session", "error", err)
		// }

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
		if err := tvwsclient.SendQuoteAddSymbolsMessage(client, qsSession, symbols); err != nil {
			slog.Error("failed to send add quote symbols session message ", "error", err)
		}
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
			// slog.Debug("data", "data.Method", data.Method)
			// slog.Debug("data", "data", data.Params)
			switch data.Method {
			case tvwsclient.MethodQuoteData:
				quoteDataMessage, err := tvwsclient.NewQuoteDataMessage(data.Params)
				if err != nil {
					slog.Error("failed to create quote data message", "error", err)
				}
				slog.Info("received quote data", "data", quoteDataMessage)
			case tvwsclient.MethodSeriesLoading:
				if len(data.Params) >= 3 {
					seriesLoadingMessage := tvwsclient.SeriesLoadingMessage{
						ChartSessionID: data.Params[0].(string),
						SeriesID:       data.Params[1].(string),
						SeriesSet:      data.Params[2].(string),
					}

					// Optional fields
					if len(data.Params) >= 4 {
						seriesLoadingMessage.SeriesNumber = data.Params[3].(string)
					}

					if len(data.Params) >= 5 {
						// Convert the interface{} back to JSON for SeriesConfig
						if configData, ok := data.Params[4].(map[string]interface{}); ok {
							if period, ok := configData["rt_update_period"].(float64); ok {
								seriesLoadingMessage.SeriesConfig = tvwsclient.SeriesConfig{
									RTUpdatePeriod: int(period),
								}
							}
						}
					}

					slog.Info("received series loading",
						"session", seriesLoadingMessage.ChartSessionID,
						"series", seriesLoadingMessage.SeriesID,
						"series set", seriesLoadingMessage.SeriesSet,
						"series number", seriesLoadingMessage.SeriesNumber,
						"config", seriesLoadingMessage.SeriesConfig,
					)
				}
			case tvwsclient.MethodSymbolResolved:
				if len(data.Params) >= 3 {
					// Convert the interface{} back to JSON for SymbolInfo
					paramJSON, err := json.Marshal(data.Params[2])
					if err != nil {
						slog.Error("failed to marshal symbol info param", "error", err)
						continue
					}

					var symbolInfo tvwsclient.SymbolInfo
					if err := json.Unmarshal(paramJSON, &symbolInfo); err != nil {
						slog.Error("failed to unmarshal symbol info", "error", err)
						continue
					}

					symbolResolvedMessage := tvwsclient.SymbolResolvedMessage{
						ChartSessionID: data.Params[0].(string),
						SeriesID:       data.Params[1].(string),
						SymbolInfo:     symbolInfo,
					}

					slog.Info("received symbol resolved",
						"session", symbolResolvedMessage.ChartSessionID,
						"series", symbolResolvedMessage.SeriesID,
						"symbol_name", symbolResolvedMessage.SymbolInfo.Name,
						"exchange", symbolResolvedMessage.SymbolInfo.Exchange,
						"description", symbolResolvedMessage.SymbolInfo.Description,
					)
				}

			case tvwsclient.MethodTimescaleUpdate:
				if len(data.Params) >= 2 {
					// Convert the interface{} back to JSON
					paramJSON, err := json.Marshal(data.Params[1])
					if err != nil {
						slog.Error("failed to marshal timescale update param", "error", err)
						continue
					}

					var timescaleUpdateData tvwsclient.TimescaleUpdateData
					if err := json.Unmarshal(paramJSON, &timescaleUpdateData); err != nil {
						slog.Error("failed to unmarshal timescale update", "error", err)
						continue
					}

					// Set the ChartSessionID from the first parameter
					timescaleUpdate := tvwsclient.TimescaleUpdateMessage{
						ChartSessionID: data.Params[0].(string),
						Data:           timescaleUpdateData,
					}

					slog.Info("received timescale update",
						"session", timescaleUpdate.ChartSessionID,
						"bar_close_time", timescaleUpdate.Data.SDS1.Lbs.BarCloseTime,
						"candles_count", len(timescaleUpdate.Data.SDS1.S),
						"last_candle", timescaleUpdate.Data.SDS1.S[len(timescaleUpdate.Data.SDS1.S)-1],
					)

					// If you want to access individual candle data:
					for _, candle := range timescaleUpdate.Data.SDS1.S {
						// candle.V contains [timestamp, open, high, low, close, volume]
						timestamp := candle.V[0]
						open := candle.V[1]
						high := candle.V[2]
						low := candle.V[3]
						close := candle.V[4]
						volume := candle.V[5]

						slog.Debug("candle data",
							"index", candle.I,
							"timestamp", timestamp,
							"open", open,
							"high", high,
							"low", low,
							"close", close,
							"volume", volume,
						)
					}
				}
			case tvwsclient.MethodSeriesCompleted:
				if len(data.Params) >= 5 {
					// Create SeriesCompletedMessage
					seriesCompletedMessage := tvwsclient.SeriesCompletedMessage{
						ChartSessionID: data.Params[0].(string),
						SeriesID:       data.Params[1].(string),
						Status:         data.Params[2].(string),
						SeriesSet:      data.Params[3].(string),
					}

					// Handle the SeriesConfig which is the 5th parameter
					if configData, ok := data.Params[4].(map[string]interface{}); ok {
						if period, ok := configData["rt_update_period"].(float64); ok {
							seriesCompletedMessage.Config = tvwsclient.SeriesConfig{
								RTUpdatePeriod: int(period),
							}
						}
					}

					slog.Info("received series completed",
						"session", seriesCompletedMessage.ChartSessionID,
						"series", seriesCompletedMessage.SeriesID,
						"status", seriesCompletedMessage.Status,
						"series_set", seriesCompletedMessage.SeriesSet,
						"config", seriesCompletedMessage.Config,
					)
				}
			case tvwsclient.MethodDataUpdate:
				if len(data.Params) >= 2 {
					// Convert the interface{} back to JSON
					paramJSON, err := json.Marshal(data.Params[1])
					if err != nil {
						slog.Error("failed to marshal du param", "error", err)
						continue
					}

					var duData tvwsclient.DuData
					if err := json.Unmarshal(paramJSON, &duData); err != nil {
						slog.Error("failed to unmarshal du data", "error", err)
						continue
					}

					duMessage := tvwsclient.DuMessage{
						ChartSessionID: data.Params[0].(string),
						Data:           duData,
					}

					// Now you can access the data like this:
					if len(duMessage.Data.SDS1.S) > 0 {
						candle := duMessage.Data.SDS1.S[0]
						slog.Info("received du update",
							"session", duMessage.ChartSessionID,
							"bar_close_time", duMessage.Data.SDS1.LBS.BarCloseTime,
							"index", candle.I,
							"timestamp", candle.V[0],
							"open", candle.V[1],
							"high", candle.V[2],
							"low", candle.V[3],
							"close", candle.V[4],
							"volume", candle.V[5],
						)
					}
				}
			}

		case <-sigChan:
			slog.Info("shutting down...")
			return
		}
	}
}
