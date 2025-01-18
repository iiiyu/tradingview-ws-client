package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dgraph-io/ristretto"
	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

type TradingViewService struct {
	client   *ent.Client
	tvClient *tvwsclient.Client
	cache    *ristretto.Cache
}
type CachedQuoteData struct {
	Name      string  `json:"name"`
	Change    float64 `json:"ch"`
	LastPrice float64 `json:"lp"`
	Timestamp int64   `json:"lp_time"`
	Volume    float64 `json:"volume"`
}

func NewTradingViewService(dbClient *ent.Client, tvClient *tvwsclient.Client, cache *ristretto.Cache) *TradingViewService {
	return &TradingViewService{
		client:   dbClient,
		tvClient: tvClient,
		cache:    cache,
	}
}

// GetDBClient returns the database client
func (s *TradingViewService) GetDBClient() *ent.Client {
	return s.client
}

// GetTVClient returns the TradingView client
func (s *TradingViewService) GetTVClient() *tvwsclient.Client {
	return s.tvClient
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
