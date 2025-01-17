package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Initialize Ent client
	client, err := ent.Open("postgres", "host=192.168.1.48 port=6543 user=postgres dbname=postgres password=uUE1yOke9wIqSAwL7bZBfKJHb5WqDnzmPIc0tlg9rF86hb5m7djpKDHulKmGy3Iy sslmode=disable")
	if err != nil {
		slog.Error("failed opening connection to postgres", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		slog.Error("failed creating schema resources", "error", err)
		os.Exit(1)
	}

	// TODO: init tradingview ws client and set read message handler
	// Initialize TradingView WebSocket client
	authToken := os.Getenv("TRADINGVIEW_AUTH_TOKEN")
	if authToken == "" {
		slog.Error("TRADINGVIEW_AUTH_TOKEN environment variable not set")
		os.Exit(1)
	}

	tvClient, err := tvwsclient.NewClient(authToken)
	if err != nil {
		slog.Error("failed to create TradingView client", "error", err)
		os.Exit(1)
	}
	defer tvClient.Close()

	// Create data channel for receiving updates
	dataChan := make(chan tvwsclient.TVResponse)

	// Start receiving data in a goroutine
	go func() {
		if err := tvClient.ReadMessage(dataChan); err != nil {
			slog.Error("failed to read messages", "error", err)
		}
	}()

	// // Start processing active sessions
	// go func() {
	// 	// Get all active sessions from the database
	// 	sessions, err := client.ActiveSession.Query().
	// 		Where(activesession.EnabledEQ(true)).
	// 		All(context.Background())
	// 	if err != nil {
	// 		slog.Error("failed to query active sessions", "error", err)
	// 		return
	// 	}

	// 	// Create quote session
	// 	qsSession := tvwsclient.GenerateSession("qs_")
	// 	if err := tvwsclient.SendQuoteCreateSessionMessage(tvClient, qsSession); err != nil {
	// 		slog.Error("failed to create quote session", "error", err)
	// 		return
	// 	}

	// 	if err := tvwsclient.SendQuoteSetFieldsMessage(tvClient, qsSession); err != nil {
	// 		slog.Error("failed to set quote fields", "error", err)
	// 		return
	// 	}

	// 	// Subscribe to each active session
	// 	for _, session := range sessions {
	// 		symbol := fmt.Sprintf("%s:%s", session.Exchange, session.Symbol)
	// 		csSession := tvwsclient.GenerateSession("cs_")

	// 		// Convert timeframe enum to TradingView format
	// 		var interval string
	// 		switch session.Timeframe {
	// 		case activesession.Timeframe10s:
	// 			interval = "10S"
	// 		case activesession.Timeframe1m:
	// 			interval = "1"
	// 		case activesession.Timeframe5m:
	// 			interval = "5"
	// 		case activesession.Timeframe1d:
	// 			interval = "D"
	// 		default:
	// 			slog.Error("invalid timeframe", "timeframe", session.Timeframe)
	// 			continue
	// 		}

	// 		if err := tvwsclient.SubscriptionChartSessionSymbol(tvClient, csSession, symbol, interval, 10); err != nil {
	// 			slog.Error("failed to subscribe to chart session",
	// 				"error", err,
	// 				"symbol", symbol,
	// 				"interval", interval)
	// 			continue
	// 		}

	// 		slog.Info("subscribed to symbol",
	// 			"symbol", symbol,
	// 			"interval", interval,
	// 			"session_id", session.SessionID)
	// 	}
	// }()

	// Process incoming data
	go func() {
		for data := range dataChan {
			switch data.Method {
			case tvwsclient.MethodQuoteData:
				quoteDataMessage, err := tvwsclient.NewQuoteDataMessage(data.Params)
				if err != nil {
					slog.Error("failed to parse quote data", "error", err)
					continue
				}
				slog.Debug("received quote data", "data", quoteDataMessage)
			case tvwsclient.MethodTimescaleUpdate:
				timescaleUpdateMessage, err := tvwsclient.NewTimescaleUpdateMessage(data.Params)
				if err != nil {
					slog.Error("failed to parse timescale update", "error", err)
					continue
				}

				// Get session ID from the message
				sessionID := timescaleUpdateMessage.ChartSessionID

				// Find the active session for this chart session
				session, err := client.ActiveSession.Query().
					Where(activesession.SessionID(sessionID)).
					Only(context.Background())
				if err != nil {
					slog.Error("failed to find active session",
						"error", err,
						"session_id", sessionID)
					continue
				}

				// Process each series update
				for _, series := range timescaleUpdateMessage.Data.SDS1.S {
					if len(series.V) < 6 {
						continue
					}

					// Extract OHLCV data
					timestamp := int64(series.V[0])
					open := series.V[1]
					high := series.V[2]
					low := series.V[3]
					close := series.V[4]
					volume := series.V[5]
					slog.Debug("series", "series.V", series.V)
					slog.Debug("received timescale update",
						"timestamp", timestamp,
						"open", open,
						"high", high,
						"low", low,
						"close", close,
						"volume", volume)

					// First check if a candle with the same key exists
					exists, err := client.Candle.Query().
						Where(
							candle.And(
								candle.ExchangeEQ(session.Exchange),
								candle.SymbolEQ(session.Symbol),
								candle.TimeframeEQ(candle.Timeframe(session.Timeframe)),
								candle.TimestampEQ(timestamp),
							),
						).Exist(context.Background())
					if err != nil {
						slog.Error("Error checking existing candle", "error", err)
						continue
					}

					if exists {
						// Update existing candle
						_, err = client.Candle.Update().
							Where(
								candle.And(
									candle.ExchangeEQ(session.Exchange),
									candle.SymbolEQ(session.Symbol),
									candle.TimeframeEQ(candle.Timeframe(session.Timeframe)),
									candle.TimestampEQ(timestamp),
								),
							).
							SetOpen(open).
							SetHigh(high).
							SetLow(low).
							SetClose(close).
							SetVolume(volume).
							Save(context.Background())
					} else {
						// Create new candle
						_, err = client.Candle.Create().
							SetExchange(session.Exchange).
							SetSymbol(session.Symbol).
							SetTimeframe(candle.Timeframe(session.Timeframe)).
							SetTimestamp(timestamp).
							SetOpen(open).
							SetHigh(high).
							SetLow(low).
							SetClose(close).
							SetVolume(volume).
							Save(context.Background())
					}

					if err != nil {
						slog.Error("Error saving candle", "error", err)
						continue
					}
				}
			case tvwsclient.MethodDataUpdate:
				duMessage, err := tvwsclient.NewDuMessage(data.Params)
				if err != nil {
					slog.Error("failed to parse data update", "error", err)
					continue
				}

				// Get session ID from the message
				sessionID := duMessage.ChartSessionID

				// Find the active session for this chart session
				session, err := client.ActiveSession.Query().
					Where(activesession.SessionID(sessionID)).
					Only(context.Background())
				if err != nil {
					slog.Error("failed to find active session",
						"error", err,
						"session_id", sessionID)
					continue
				}

				// Process each series update
				for _, series := range duMessage.Data.SDS1.S {
					if len(series.V) < 6 {
						continue
					}

					// Extract OHLCV data
					timestamp := int64(series.V[0])
					open := series.V[1]
					high := series.V[2]
					low := series.V[3]
					close := series.V[4]
					volume := series.V[5]

					// First check if a candle with the same key exists
					exists, err := client.Candle.Query().
						Where(
							candle.And(
								candle.ExchangeEQ(session.Exchange),
								candle.SymbolEQ(session.Symbol),
								candle.TimeframeEQ(candle.Timeframe(session.Timeframe)),
								candle.TimestampEQ(timestamp),
							),
						).Exist(context.Background())
					if err != nil {
						slog.Error("Error checking existing candle", "error", err)
						continue
					}

					if exists {
						// Update existing candle
						_, err = client.Candle.Update().
							Where(
								candle.And(
									candle.ExchangeEQ(session.Exchange),
									candle.SymbolEQ(session.Symbol),
									candle.TimeframeEQ(candle.Timeframe(session.Timeframe)),
									candle.TimestampEQ(timestamp),
								),
							).
							SetOpen(open).
							SetHigh(high).
							SetLow(low).
							SetClose(close).
							SetVolume(volume).
							Save(context.Background())
					} else {
						// Create new candle
						_, err = client.Candle.Create().
							SetExchange(session.Exchange).
							SetSymbol(session.Symbol).
							SetTimeframe(candle.Timeframe(session.Timeframe)).
							SetTimestamp(timestamp).
							SetOpen(open).
							SetHigh(high).
							SetLow(low).
							SetClose(close).
							SetVolume(volume).
							Save(context.Background())
					}

					if err != nil {
						slog.Error("Error saving candle", "error", err)
						continue
					}
				}
			}
		}
	}()

	app := fiber.New(fiber.Config{
		AppName: "TradingView Data Service",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			slog.Error("fiber error",
				"error", err,
				"path", c.Path(),
				"method", c.Method(),
				"ip", c.IP(),
			)

			// Default 500 status code
			code := fiber.StatusInternalServerError

			// Check if it's a fiber.Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Basic routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "TradingView Data Service",
			"status":  "running",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Symbol management routes
	app.Post("/symbols", func(c *fiber.Ctx) error {
		var input struct {
			SessionID string `json:"session_id"`
			Exchange  string `json:"exchange"`
			Symbol    string `json:"symbol"`
			Timeframe string `json:"timeframe"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// Create chart session
		csSession := tvwsclient.GenerateSession("cs_")
		symbol := fmt.Sprintf("%s:%s", input.Exchange, input.Symbol)

		// Convert timeframe to TradingView format
		var interval string
		switch activesession.Timeframe(input.Timeframe) {
		case activesession.Timeframe10S:
			interval = "10S"
		case activesession.Timeframe1:
			interval = "1"
		case activesession.Timeframe5:
			interval = "5"
		case activesession.Timeframe1D:
			interval = "1D"
		default:
			return c.Status(400).JSON(fiber.Map{"error": "invalid timeframe"})
		}

		// Subscribe to TradingView
		if err := tvwsclient.SubscriptionChartSessionSymbol(tvClient, csSession, symbol, interval, 10); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to subscribe to TradingView: " + err.Error()})
		}

		// Save to database
		session, err := client.ActiveSession.Create().
			SetSessionID(csSession). // Use the generated chart session ID
			SetExchange(input.Exchange).
			SetSymbol(input.Symbol).
			SetTimeframe(activesession.Timeframe(input.Timeframe)).
			SetEnabled(true).
			Save(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(session)
	})

	app.Delete("/symbols/:session_id", func(c *fiber.Ctx) error {
		sessionID := c.Params("session_id")
		_, err := client.ActiveSession.Delete().
			Where(activesession.SessionID(sessionID)).
			Exec(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(204)
	})

	app.Get("/symbols", func(c *fiber.Ctx) error {
		sessions, err := client.ActiveSession.Query().
			Where(activesession.EnabledEQ(true)).
			All(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(sessions)
	})

	app.Get("/symbols/:exchange/:symbol", func(c *fiber.Ctx) error {
		exchange := c.Params("exchange")
		symbol := c.Params("symbol")

		sessions, err := client.ActiveSession.Query().
			Where(
				activesession.EnabledEQ(true),
				activesession.ExchangeEQ(exchange),
				activesession.SymbolEQ(symbol),
			).
			All(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(sessions)
	})

	app.Get("/symbols/session/:session_id/status", func(c *fiber.Ctx) error {
		sessionID := c.Params("session_id")
		session, err := client.ActiveSession.Query().
			Where(activesession.SessionID(sessionID)).
			Only(c.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				return c.Status(404).JSON(fiber.Map{"error": "Session not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(session)
	})

	// Candlestick data routes
	app.Get("/candles/:exchange/:symbol", func(c *fiber.Ctx) error {
		exchange := c.Params("exchange")
		symbol := c.Params("symbol")
		timeframe := c.Query("timeframe", "1")
		limit := 100 // Default limit

		candles, err := client.Candle.Query().
			Where(
				candle.ExchangeEQ(exchange),
				candle.SymbolEQ(symbol),
				candle.TimeframeEQ(candle.Timeframe(timeframe)),
			).
			Order(ent.Desc("timestamp")).
			Limit(limit).
			All(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(candles)
	})

	slog.Info("starting server on port 3333")
	if err := app.Listen(":3333"); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
