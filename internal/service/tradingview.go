package service

import (
	"context"
	"log/slog"

	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/pkg/tvwsclient"
)

type TradingViewService struct {
	client   *ent.Client
	tvClient *tvwsclient.Client
}

func NewTradingViewService(dbClient *ent.Client, tvClient *tvwsclient.Client) *TradingViewService {
	return &TradingViewService{
		client:   dbClient,
		tvClient: tvClient,
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
				candle.TimeframeEQ(candle.Timeframe(session.Timeframe)),
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
		_, err = s.client.Candle.Create().
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

	return err
}
