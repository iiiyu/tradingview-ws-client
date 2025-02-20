package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dgraph-io/ristretto"
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

	// Initialize Ent entClient
	entClient, err := ent.Open("postgres", cfg.GetDSN())
	if err != nil {
		slog.Error("failed opening connection to postgres", "error", err)
		os.Exit(1)
	}
	defer entClient.Close()

	// Run the auto migration tool
	if err := entClient.Schema.Create(context.Background()); err != nil {
		slog.Error("failed creating schema resources", "error", err)
		os.Exit(1)
	}

	// Clean up old sessions
	if err := service.CleanUpOldSessions(entClient); err != nil {
		slog.Error("failed to clean up old sessions", "error", err)
		os.Exit(1)
	}

	// Initialize AuthTokenManager
	deviceToken, sessionID, sessionSign := cfg.GetTradingViewConfig()
	tvwsclient.InitAuthTokenManager(tvwsclient.NewTVHttpClient("https://www.tradingview.com", deviceToken, sessionID, sessionSign))

	// Initialize TradingView WebSocket client
	tvClient, err := tvwsclient.NewClient()
	if err != nil {
		slog.Error("failed to create TradingView client", "error", err)
		os.Exit(1)
	}
	defer tvClient.Close()

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		slog.Error("failed to create cache", "error", err)
		os.Exit(1)
	}

	// Initialize service
	tvService := service.NewTradingViewService(entClient, tvClient, cache)

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
	port := cfg.Port
	slog.Info("starting server on port", "port", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
