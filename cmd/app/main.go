package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
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

		session, err := client.ActiveSession.Create().
			SetSessionID(input.SessionID).
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
