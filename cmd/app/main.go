package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"

	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/internal/config"
	"github.com/iiiyu/tradingview-ws-client/internal/handler"
	"github.com/iiiyu/tradingview-ws-client/internal/service"
	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize Ent client
	client, err := ent.Open("postgres", cfg.GetDSN())
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

	// Initialize TradingView WebSocket client
	tvClient, err := tvwsclient.NewClient(cfg.TradingView.AuthToken)
	if err != nil {
		slog.Error("failed to create TradingView client", "error", err)
		os.Exit(1)
	}
	defer tvClient.Close()

	// Initialize service
	tvService := service.NewTradingViewService(client, tvClient)

	// Create data channel for receiving updates
	dataChan := make(chan tvwsclient.TVResponse)

	// Start receiving data in a goroutine
	go func() {
		if err := tvClient.ReadMessage(dataChan); err != nil {
			slog.Error("failed to read messages", "error", err)
		}
	}()

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

				if err := tvService.ProcessTimescaleUpdate(timescaleUpdateMessage); err != nil {
					slog.Error("failed to process timescale update", "error", err)
				}

			case tvwsclient.MethodDataUpdate:
				duMessage, err := tvwsclient.NewDuMessage(data.Params)
				if err != nil {
					slog.Error("failed to parse data update", "error", err)
					continue
				}

				if err := tvService.ProcessDataUpdate(duMessage); err != nil {
					slog.Error("failed to process data update", "error", err)
				}
			}
		}
	}()

	// Initialize Fiber app and handlers
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

	h := handler.NewHandler(tvService)
	h.RegisterRoutes(app)

	// Start the server
	slog.Info("starting server on port 3333")
	if err := app.Listen(":3333"); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
