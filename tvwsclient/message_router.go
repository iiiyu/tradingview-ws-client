package tvwsclient

import (
	"context"
	"fmt"
	"log/slog"
)

// MessageRouter handles routing different message types to appropriate handlers
type MessageRouter struct {
	handlers map[string]MessageHandler
	logger   *slog.Logger
}

// NewMessageRouter creates a new message router
func NewMessageRouter(logger *slog.Logger) *MessageRouter {
	return &MessageRouter{
		handlers: make(map[string]MessageHandler),
		logger:   logger,
	}
}

// RegisterHandler registers a handler for a specific message method
func (r *MessageRouter) RegisterHandler(method string, handler MessageHandler) {
	r.handlers[method] = handler
}

// RouteMessage routes a message to the appropriate handler
func (r *MessageRouter) RouteMessage(ctx context.Context, response TVResponse) error {
	handler, exists := r.handlers[response.Method]
	if !exists {
		r.logger.Debug("no handler registered for message method", "method", response.Method)
		return nil // Not an error, just no handler
	}

	switch response.Method {
	case MethodQuoteData:
		msg, err := NewQuoteDataMessage(response.Params)
		if err != nil {
			return WrapMessageError("route.quote_data", err)
		}
		return handler.HandleQuoteData(ctx, msg)

	case MethodTimescaleUpdate:
		msg, err := NewTimescaleUpdateMessage(response.Params)
		if err != nil {
			return WrapMessageError("route.timescale_update", err)
		}
		return handler.HandleTimescaleUpdate(ctx, msg)

	case MethodDataUpdate:
		msg, err := NewDuMessage(response.Params)
		if err != nil {
			return WrapMessageError("route.data_update", err)
		}
		return handler.HandleDataUpdate(ctx, msg)

	case MethodQuoteCompleted:
		msg, err := NewQuoteCompletedMessage(response.Params)
		if err != nil {
			return WrapMessageError("route.quote_completed", err)
		}
		return handler.HandleQuoteCompleted(ctx, msg)

	default:
		r.logger.Debug("unhandled message method", "method", response.Method)
		return nil
	}
}

// DefaultMessageHandler provides default implementations for message handling
type DefaultMessageHandler struct {
	quoteProcessor     QuoteProcessor
	candleProcessor    CandleProcessor
	sessionProcessor   SessionProcessor
	logger             *slog.Logger
}

// QuoteProcessor defines interface for processing quote data
type QuoteProcessor interface {
	ProcessQuoteData(ctx context.Context, msg *QuoteDataMessage) error
}

// CandleProcessor defines interface for processing candle data
type CandleProcessor interface {
	ProcessTimescaleUpdate(ctx context.Context, msg *TimescaleUpdateMessage) error
	ProcessDataUpdate(ctx context.Context, msg *DuMessage) error
}

// SessionProcessor defines interface for processing session events
type SessionProcessor interface {
	ProcessQuoteCompleted(ctx context.Context, msg *QuoteCompletedMessage) error
}

// NewDefaultMessageHandler creates a new default message handler
func NewDefaultMessageHandler(
	quoteProc QuoteProcessor,
	candleProc CandleProcessor,
	sessionProc SessionProcessor,
	logger *slog.Logger,
) *DefaultMessageHandler {
	return &DefaultMessageHandler{
		quoteProcessor:   quoteProc,
		candleProcessor:  candleProc,
		sessionProcessor: sessionProc,
		logger:           logger,
	}
}

// HandleQuoteData implements MessageHandler
func (h *DefaultMessageHandler) HandleQuoteData(ctx context.Context, msg *QuoteDataMessage) error {
	if h.quoteProcessor == nil {
		return fmt.Errorf("quote processor not configured")
	}
	return h.quoteProcessor.ProcessQuoteData(ctx, msg)
}

// HandleTimescaleUpdate implements MessageHandler
func (h *DefaultMessageHandler) HandleTimescaleUpdate(ctx context.Context, msg *TimescaleUpdateMessage) error {
	if h.candleProcessor == nil {
		return fmt.Errorf("candle processor not configured")
	}
	return h.candleProcessor.ProcessTimescaleUpdate(ctx, msg)
}

// HandleDataUpdate implements MessageHandler
func (h *DefaultMessageHandler) HandleDataUpdate(ctx context.Context, msg *DuMessage) error {
	if h.candleProcessor == nil {
		return fmt.Errorf("candle processor not configured")
	}
	return h.candleProcessor.ProcessDataUpdate(ctx, msg)
}

// HandleQuoteCompleted implements MessageHandler
func (h *DefaultMessageHandler) HandleQuoteCompleted(ctx context.Context, msg *QuoteCompletedMessage) error {
	if h.sessionProcessor == nil {
		h.logger.Debug("session processor not configured, ignoring quote completed")
		return nil
	}
	return h.sessionProcessor.ProcessQuoteCompleted(ctx, msg)
}

// MessageProcessor combines all processors for convenience
type MessageProcessor struct {
	*DefaultMessageHandler
}

// NewMessageProcessor creates a new message processor with all handlers
func NewMessageProcessor(
	quoteProc QuoteProcessor,
	candleProc CandleProcessor,
	sessionProc SessionProcessor,
	logger *slog.Logger,
) *MessageProcessor {
	return &MessageProcessor{
		DefaultMessageHandler: NewDefaultMessageHandler(quoteProc, candleProc, sessionProc, logger),
	}
}