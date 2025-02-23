package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

type TradingViewService struct {
	client     *ent.Client
	tvClient   *tvwsclient.Client
	cache      *ristretto.Cache
	mu         sync.Mutex
	cancelFunc context.CancelFunc
}
type CachedQuoteData struct {
	Name      string  `json:"name"`
	Change    float64 `json:"ch"`
	LastPrice float64 `json:"lp"`
	Timestamp int64   `json:"lp_time"`
	Volume    float64 `json:"volume"`
	Bid       float64 `json:"bid"`
	Ask       float64 `json:"ask"`
	BidSize   int     `json:"bid_size"`
	AskSize   int     `json:"ask_size"`
	// RCH (Regular Change): The absolute price change during regular trading hours
	RCH float64 `json:"rch,omitempty"`
	// RCHP (Regular Change Percentage): The percentage change during regular trading hours
	RCHP float64 `json:"rchp,omitempty"`
	// RTC (Real-Time Close): The current/latest closing price in real-time
	RTC float64 `json:"rtc,omitempty"`
	// RTC_Time: The timestamp of the latest real-time close price
	RTC_Time int64 `json:"rtc_time,omitempty"`
}

func NewTradingViewService(dbClient *ent.Client, tvClient *tvwsclient.Client, cache *ristretto.Cache) *TradingViewService {
	service := &TradingViewService{
		client:   dbClient,
		tvClient: tvClient,
		cache:    cache,
	}

	service.readTradingViewMessage()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			if tvwsclient.GetAuthTokenManager().CheckAuthTokenExpired() {
				slog.Info("auth token expired, reconnecting TradingView client")
				if err := service.ReconnectTVClient(); err != nil {
					slog.Error("failed to reconnect TradingView client", "error", err)
				}
			}
		}
	}()

	return service
}

func (s *TradingViewService) readTradingViewMessage() {
	s.mu.Lock()
	// Cancel previous goroutines if they exist
	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	// Create new context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.mu.Unlock()

	// Create data channel for receiving updates
	dataChan := make(chan tvwsclient.TVResponse)

	// Start receiving data in a goroutine
	go func() {
		defer close(dataChan)
		if err := s.tvClient.ReadMessage(dataChan); err != nil {
			slog.Error("failed to read messages", "error", err)
		}
	}()

	// Process incoming data
	go func() {
		defer slog.Info("message processing goroutine stopped")
		for {
			select {
			case <-ctx.Done():
				slog.Error("ctx.Done()")
				return
			case data, ok := <-dataChan:
				if !ok {
					slog.Error("dataChan closed", "data", data, "ok", ok)
					return
				}
				slog.Debug("Received message", "method", data)
				switch data.Method {
				case tvwsclient.MethodQuoteCompleted:
					slog.Debug("MethodQuoteCompleted", "data", data)
					// quoteCompletedMessage, err := tvwsclient.NewQuoteCompletedMessage(data.Params)
					// if err != nil {
					// 	slog.Error("failed to parse quote completed", "error", err)
					// 	continue
					// }

					// if err := tvwsclient.SendQuoteCompletedMessageAfterQuoteCompleted(s.tvClient, quoteCompletedMessage.SessionID, quoteCompletedMessage.ReceivedMessage); err != nil {
					// 	slog.Error("failed to send quote completed message after quote completed", "error", err)
					// }

					// session := tvwsclient.GenerateSession("qs_")
					// if err := tvwsclient.SubscriptionQuoteSessionSymbol(s.tvClient, session, quoteCompletedMessage.ReceivedMessage); err != nil {
					// 	slog.Error("failed to subscribe to quote session symbol", "error", err)
					// }

				case tvwsclient.MethodQuoteData:
					quoteDataMessage, err := tvwsclient.NewQuoteDataMessage(data.Params)
					if err != nil {
						slog.Error("failed to parse quote data", "error", err)
						continue
					}

					if err := s.ProcessQuoteData(quoteDataMessage); err != nil {
						slog.Error("failed to process quote data", "error", err)
					}

				case tvwsclient.MethodTimescaleUpdate:
					timescaleUpdateMessage, err := tvwsclient.NewTimescaleUpdateMessage(data.Params)
					if err != nil {
						slog.Error("failed to parse timescale update", "error", err)
						continue
					}

					if err := s.ProcessTimescaleUpdate(timescaleUpdateMessage); err != nil {
						slog.Error("failed to process timescale update", "error", err)
					}

				case tvwsclient.MethodDataUpdate:
					duMessage, err := tvwsclient.NewDuMessage(data.Params)
					if err != nil {
						slog.Error("failed to parse data update", "error", err)
						continue
					}

					if err := s.ProcessDataUpdate(duMessage); err != nil {
						slog.Error("failed to process data update", "error", err)
					}
				}
			}
		}
	}()
}

// GetDBClient returns the database client
func (s *TradingViewService) GetDBClient() *ent.Client {
	return s.client
}

// GetTVClient returns the TradingView client
func (s *TradingViewService) GetTVClient() *tvwsclient.Client {
	return s.tvClient
}

func (s *TradingViewService) Unsubscribe(sessions []*ent.ActiveSession) error {

	var err error
	// Unsubscribe each session based on its type
	for _, session := range sessions {
		// Unsubscribe from TradingView based on session type
		if session.Type == activesession.TypeCandle {
			err = tvwsclient.SendChartDeleteSessionMessage(s.GetTVClient(), session.SessionID)
		} else if session.Type == activesession.TypeQuote {
			err = tvwsclient.SendQuoteRemoveSymbolsMessage(s.GetTVClient(), session.SessionID,
				[]string{fmt.Sprintf("%s:%s", session.Exchange, session.Symbol)})
		}

		if err != nil {
			// ignore error if it fails
			slog.Error("failed to unsubscribe from TradingView", "error", err)
		}

		// Update session status to disabled
		_, err = session.Update().
			SetEnabled(false).
			Save(context.Background())

		if err != nil {
			// ignore error if it fails
			slog.Error("failed to update session status", "error", err)
		}
	}

	return nil
}

func (s *TradingViewService) Subscribe(sessions []*ent.ActiveSession) error {

	for _, session := range sessions {
		symbol := fmt.Sprintf("%s:%s", session.Exchange, session.Symbol)
		var sessionID string

		switch session.Type {
		case activesession.TypeCandle:
			sessionID = tvwsclient.GenerateSession("cs_")

			// Subscribe to TradingView
			if err := tvwsclient.SubscriptionChartSessionSymbol(s.tvClient, sessionID, symbol, session.Timeframe.String(), 300); err != nil {
				return fmt.Errorf("failed to subscribe to TradingView for session %s: %w", session.SessionID, err)
			}

		case activesession.TypeQuote:
			sessionID = tvwsclient.GenerateSession("qs_")
			if err := tvwsclient.SubscriptionQuoteSessionSymbol(s.tvClient, sessionID, symbol); err != nil {
				return fmt.Errorf("failed to subscribe to TradingView for session %s: %w", session.SessionID, err)
			}

		default:
			return fmt.Errorf("invalid session type for session: %s", session.SessionID)
		}

		// Update session to enabled
		_, err := session.Update().
			SetEnabled(true).
			SetSessionID(sessionID).
			Save(context.Background())

		if err != nil {
			return fmt.Errorf("failed to update session %s status: %w", session.SessionID, err)
		}
	}
	return nil
}

// ReconnectTVClient reconnects the TradingView client
func (s *TradingViewService) ReconnectTVClient() error {
	ctx := context.Background()

	// 1. fetch all active sessions and keep the symbol name
	activeSessions, err := s.client.ActiveSession.Query().
		Where(activesession.EnabledEQ(true)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch active sessions: %w", err)
	}

	// 2. try to unsubscribe from all symbols, ignore error if it fails
	err = s.Unsubscribe(activeSessions)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from TradingView: %w", err)
	}

	// 3. reconnect to TradingView
	if err := s.tvClient.Reconnect(); err != nil {
		return fmt.Errorf("failed to reconnect: %w", err)
	}

	s.readTradingViewMessage()

	// 4. subscribe to all symbols and reactivate sessions
	err = s.Subscribe(activeSessions)
	if err != nil {
		return fmt.Errorf("failed to subscribe to TradingView: %w", err)
	}

	return nil
}

func (s *TradingViewService) ProcessTimescaleUpdate(msg *tvwsclient.TimescaleUpdateMessage) error {
	// Get session ID from the message
	sessionID := msg.ChartSessionID

	// Find the active session for this chart session
	session, err := s.client.ActiveSession.Query().
		Where(activesession.SessionID(sessionID)).
		Only(context.Background())
	if err != nil {
		slog.Error("failed to find active session",
			"error", err,
			"session_id", sessionID)
		return err
	}

	// Process each series update
	for _, series := range msg.Data.SDS1.S {
		if len(series.V) < 6 {
			continue
		}

		if err := s.processCandleData(session, series.V); err != nil {
			slog.Error("failed to process candle", "error", err)
			continue
		}
	}
	return nil
}

func (s *TradingViewService) ProcessDataUpdate(msg *tvwsclient.DuMessage) error {
	sessionID := msg.ChartSessionID

	session, err := s.client.ActiveSession.Query().
		Where(activesession.SessionID(sessionID)).
		Only(context.Background())
	if err != nil {
		slog.Error("failed to find active session",
			"error", err,
			"session_id", sessionID)
		return err
	}

	for _, series := range msg.Data.SDS1.S {
		if len(series.V) < 6 {
			continue
		}

		if err := s.processCandleData(session, series.V); err != nil {
			slog.Error("failed to process candle", "error", err)
			continue
		}
	}
	return nil
}

func (s *TradingViewService) ProcessQuoteData(msg *tvwsclient.QuoteDataMessage) error {
	// slog.Debug("ProcessQuoteData", "msg", msg)
	// Extract symbol name from the message
	symbolName := msg.Data.Name

	// Try to get existing cached data
	var cachedData *CachedQuoteData
	if existingValue, found := s.cache.Get(symbolName); found {
		if existing, ok := existingValue.(*CachedQuoteData); ok {
			cachedData = existing
		}
	}

	// If no existing data, create new
	if cachedData == nil {
		cachedData = &CachedQuoteData{Name: symbolName}
	}

	// Update only non-zero values
	if msg.Data.Values.Change != 0 {
		cachedData.Change = msg.Data.Values.Change
	}
	if msg.Data.Values.LastPrice != 0 {
		cachedData.LastPrice = msg.Data.Values.LastPrice
	}
	if msg.Data.Values.LastPriceTime != 0 {
		cachedData.Timestamp = msg.Data.Values.LastPriceTime
	}
	if msg.Data.Values.Volume != 0 {
		cachedData.Volume = msg.Data.Values.Volume
	}
	if msg.Data.Values.Bid != 0 {
		cachedData.Bid = msg.Data.Values.Bid
	}
	if msg.Data.Values.Ask != 0 {
		cachedData.Ask = msg.Data.Values.Ask
	}
	if msg.Data.Values.BidSize != 0 {
		cachedData.BidSize = msg.Data.Values.BidSize
	}
	if msg.Data.Values.AskSize != 0 {
		cachedData.AskSize = msg.Data.Values.AskSize
	}

	// real-time close price
	if msg.Data.Values.RCH != 0 {
		cachedData.RCH = msg.Data.Values.RCH
	}
	if msg.Data.Values.RCHP != 0 {
		cachedData.RCHP = msg.Data.Values.RCHP
	}
	if msg.Data.Values.RTC != 0 {
		cachedData.RTC = msg.Data.Values.RTC
	}
	if msg.Data.Values.RTC_Time != 0 {
		cachedData.RTC_Time = msg.Data.Values.RTC_Time
	}

	// Set the data in cache
	// The third parameter (1) is the cost of storing this item
	if !s.cache.Set(symbolName, cachedData, 1) {
		slog.Warn("failed to set quote data in cache",
			"symbol", symbolName,
			"data", cachedData)
	}

	// Wait for cache write
	s.cache.Wait()

	return nil
}

func (s *TradingViewService) processCandleData(session *ent.ActiveSession, data []float64) error {
	timestamp := int64(data[0])
	open := data[1]
	high := data[2]
	low := data[3]
	close := data[4]
	volume := data[5]

	exists, err := s.client.Candle.Query().
		Where(
			candle.And(
				candle.ExchangeEQ(session.Exchange),
				candle.SymbolEQ(session.Symbol),
				candle.TimeframeEQ(candle.Timeframe(string(*session.Timeframe))),
				candle.TimestampEQ(timestamp),
			),
		).Exist(context.Background())
	if err != nil {
		return err
	}

	if exists {
		_, err = s.client.Candle.Update().
			Where(
				candle.And(
					candle.ExchangeEQ(session.Exchange),
					candle.SymbolEQ(session.Symbol),
					candle.TimeframeEQ(candle.Timeframe(string(*session.Timeframe))),
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
		_, err = s.client.Candle.Create().
			SetExchange(session.Exchange).
			SetSymbol(session.Symbol).
			SetTimeframe(candle.Timeframe(string(*session.Timeframe))).
			SetTimestamp(timestamp).
			SetOpen(open).
			SetHigh(high).
			SetLow(low).
			SetClose(close).
			SetVolume(volume).
			Save(context.Background())
	}

	return err
}

func (s *TradingViewService) GetQuoteData(symbol string) (*CachedQuoteData, error) {
	value, found := s.cache.Get(symbol)
	if !found {
		return nil, fmt.Errorf("quote data not found for symbol: %s", symbol)
	}

	data, ok := value.(*CachedQuoteData)
	if !ok {
		return nil, fmt.Errorf("invalid cache data type for symbol: %s", symbol)
	}

	return data, nil
}
