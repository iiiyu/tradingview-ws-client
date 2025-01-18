package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/ent/predicate"
	"github.com/iiiyu/tradingview-ws-client/internal/service"
	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

type Handler struct {
	tvService *service.TradingViewService
}

type SymbolRequest struct {
	Exchange  string `json:"exchange"`
	Symbol    string `json:"symbol"`
	Timeframe string `json:"timeframe"`
	Type      string `json:"type"`
}

func NewHandler(tvService *service.TradingViewService) *Handler {
	return &Handler{
		tvService: tvService,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Basic routes
	app.Get("/", h.handleHome)
	app.Get("/health", h.handleHealth)

	// Symbol management routes
	// add symbol
	app.Post("/symbols", h.handleCreateSymbol)
	// delete symbol
	app.Delete("/symbols", h.handleDeleteSymbol)
	// unsubscribe all symbols
	app.Delete("/symbols/unsubscribe", h.handleUnsubscribeAllSymbols)
	// list symbols
	app.Get("/symbols", h.handleListSymbols)
	// get symbol by exchange
	app.Get("/symbols/:exchange/:symbol", h.handleGetSymbolByExchange)
	// get symbol status
	app.Get("/symbols/session/:session_id/status", h.handleGetSymbolStatus)

	// Candlestick data routes
	app.Get("/candles/:exchange/:symbol/:timeframe/:limit", h.handleGetCandles)

	// Quote data routes
	app.Get("/quotes/:exchange/:symbol", h.handleGetQuoteData)
}

func (h *Handler) handleHome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"service": "TradingView Data Service",
		"status":  "running",
	})
}

func (h *Handler) handleHealth(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func (h *Handler) handleCreateSymbol(c *fiber.Ctx) error {
	input := SymbolRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// check if session exists
	// Build query conditions
	conditions := []predicate.ActiveSession{
		activesession.ExchangeEQ(input.Exchange),
		activesession.SymbolEQ(input.Symbol),
		activesession.TypeEQ(activesession.Type(input.Type)),
		activesession.EnabledEQ(true),
	}

	// Add timeframe condition for candle type
	if activesession.Type(input.Type) == activesession.TypeCandle {
		conditions = append(conditions, activesession.TimeframeEQ(activesession.Timeframe(input.Timeframe)))
	}

	// Find active session
	session, err := h.tvService.GetDBClient().ActiveSession.Query().
		Where(conditions...).
		Only(c.Context())

	if err != nil {
		if !ent.IsNotFound(err) {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if session != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Session already exists"})
	}

	switch activesession.Type(input.Type) {
	case activesession.TypeCandle:
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
		if err := tvwsclient.SubscriptionChartSessionSymbol(h.tvService.GetTVClient(), csSession, symbol, interval, 10); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to subscribe to TradingView: " + err.Error()})
		}

		// Save to database
		session, err := h.tvService.GetDBClient().ActiveSession.Create().
			SetSessionID(csSession).
			SetExchange(input.Exchange).
			SetSymbol(input.Symbol).
			SetTimeframe(activesession.Timeframe(input.Timeframe)).
			SetType(activesession.Type(input.Type)).
			SetEnabled(true).
			Save(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(session)
	case activesession.TypeQuote:
		// Create quote session
		qsSession := tvwsclient.GenerateSession("qs_")
		symbol := fmt.Sprintf("%s:%s", input.Exchange, input.Symbol)

		// Subscribe to TradingView
		if err := tvwsclient.SubscriptionQuoteSessionSymbol(h.tvService.GetTVClient(), qsSession, symbol); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to subscribe to TradingView: " + err.Error()})
		}

		// Save to database
		session, err := h.tvService.GetDBClient().ActiveSession.Create().
			SetSessionID(qsSession).
			SetExchange(input.Exchange).
			SetSymbol(input.Symbol).
			SetType(activesession.Type(input.Type)).
			SetEnabled(true).
			Save(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(session)
	}
	return c.JSON(fiber.Map{"error": "invalid type"})
}

func (h *Handler) handleDeleteSymbol(c *fiber.Ctx) error {
	input := SymbolRequest{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// 1. check if session exists
	// use type ,exchange,symbol and if type = candle add timeframe to check if session exists

	// Build query conditions
	conditions := []predicate.ActiveSession{
		activesession.ExchangeEQ(input.Exchange),
		activesession.SymbolEQ(input.Symbol),
		activesession.TypeEQ(activesession.Type(input.Type)),
		activesession.EnabledEQ(true),
	}

	// Add timeframe condition for candle type
	if activesession.Type(input.Type) == activesession.TypeCandle {
		conditions = append(conditions, activesession.TimeframeEQ(activesession.Timeframe(input.Timeframe)))
	}

	// Find active session
	session, err := h.tvService.GetDBClient().ActiveSession.Query().
		Where(conditions...).
		Only(c.Context())

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(404).JSON(fiber.Map{"error": "Session not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. use sessionID to unsubscribe from TradingView
	// Unsubscribe from TradingView based on session type
	if session.Type == activesession.TypeCandle {
		err = tvwsclient.SendChartDeleteSessionMessage(h.tvService.GetTVClient(), session.SessionID)
	} else if session.Type == activesession.TypeQuote {
		err = tvwsclient.SendQuoteRemoveSymbolsMessage(h.tvService.GetTVClient(), session.SessionID, []string{fmt.Sprintf("%s:%s", session.Exchange, session.Symbol)})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to unsubscribe from TradingView: " + err.Error()})
	}

	// 3. update session enabled to false
	// Update session status
	session, err = session.Update().
		SetEnabled(false).
		Save(c.Context())

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update session status: " + err.Error()})
	}

	return c.JSON(session)
}

func (h *Handler) handleUnsubscribeAllSymbols(c *fiber.Ctx) error {
	// Get all enabled sessions
	sessions, err := h.tvService.GetDBClient().ActiveSession.Query().
		Where(activesession.EnabledEQ(true)).
		All(c.Context())

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch active sessions: " + err.Error()})
	}

	// Unsubscribe each session based on its type
	for _, session := range sessions {
		// Unsubscribe from TradingView based on session type
		if session.Type == activesession.TypeCandle {
			err = tvwsclient.SendChartDeleteSessionMessage(h.tvService.GetTVClient(), session.SessionID)
		} else if session.Type == activesession.TypeQuote {
			err = tvwsclient.SendQuoteRemoveSymbolsMessage(h.tvService.GetTVClient(), session.SessionID,
				[]string{fmt.Sprintf("%s:%s", session.Exchange, session.Symbol)})
		}

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to unsubscribe from TradingView: " + err.Error()})
		}

		// Update session status to disabled
		_, err = session.Update().
			SetEnabled(false).
			Save(c.Context())

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update session status: " + err.Error()})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Successfully unsubscribed all symbols",
		"count":   len(sessions),
	})
}

func (h *Handler) handleListSymbols(c *fiber.Ctx) error {
	sessions, err := h.tvService.GetDBClient().ActiveSession.Query().
		Where(activesession.EnabledEQ(true)).
		All(c.Context())

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(sessions)
}

func (h *Handler) handleGetSymbolByExchange(c *fiber.Ctx) error {
	exchange := c.Params("exchange")
	symbol := c.Params("symbol")

	sessions, err := h.tvService.GetDBClient().ActiveSession.Query().
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
}

func (h *Handler) handleGetSymbolStatus(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")
	session, err := h.tvService.GetDBClient().ActiveSession.Query().
		Where(activesession.SessionID(sessionID)).
		Only(c.Context())

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(404).JSON(fiber.Map{"error": "Session not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(session)
}

func (h *Handler) handleGetCandles(c *fiber.Ctx) error {
	exchange := c.Params("exchange")
	symbol := c.Params("symbol")
	timeframe := c.Params("timeframe", "1")
	limitStr := c.Params("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid limit parameter"})
	}

	candles, err := h.tvService.GetDBClient().Candle.Query().
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
}

func (h *Handler) handleGetQuoteData(c *fiber.Ctx) error {
	exchange := c.Params("exchange")
	symbol := c.Params("symbol")
	querySymbol := fmt.Sprintf("%s:%s", exchange, symbol)

	quoteData, err := h.tvService.GetQuoteData(querySymbol)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(quoteData)
}
