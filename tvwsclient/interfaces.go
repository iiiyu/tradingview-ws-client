package tvwsclient

import (
	"context"
)

// TradingViewClient defines the interface for TradingView WebSocket client operations
type TradingViewClient interface {
	// Connection management
	Close() error
	Reconnect() error
	
	// Message reading
	ReadMessage(dataChan chan<- TVResponse) error
	
	// Initialization
	SendInitMessage() error
}

// AuthTokenManagerInterface defines the interface for authentication token management
type AuthTokenManagerInterface interface {
	GetToken() string
	CheckAuthTokenExpired() bool
	RefreshToken() error
}

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	DoRequest(method, url string, body []byte) ([]byte, error)
}

// MessageSender defines the interface for sending WebSocket messages
type MessageSender interface {
	// Auth messages
	SendSetAuthTokenMessage(authToken string) error
	SendSetLocalMessage() error
	
	// Chart messages
	SendChartCreateSessionMessage(session string) error
	SendChartDeleteSessionMessage(session string) error
	SendSwitchTimezoneMessage(session string) error
	SendResolveSymbolMessage(session, symbol string) error
	SendCreateSeriesMessage(session, interval string, seriesNumber int64) error
	
	// Quote messages
	SendQuoteCreateSessionMessage(session string) error
	SendQuoteSetFieldsMessage(session string) error
	SendQuoteRemoveSymbolsMessage(session string, symbols []string) error
}

// MessageHandler defines the interface for handling different message types
type MessageHandler interface {
	HandleQuoteData(ctx context.Context, msg *QuoteDataMessage) error
	HandleTimescaleUpdate(ctx context.Context, msg *TimescaleUpdateMessage) error
	HandleDataUpdate(ctx context.Context, msg *DuMessage) error
	HandleQuoteCompleted(ctx context.Context, msg *QuoteCompletedMessage) error
}

// Repository defines the interface for data persistence operations
type Repository interface {
	// ActiveSession operations
	CreateActiveSession(ctx context.Context, session *ActiveSessionData) error
	UpdateActiveSession(ctx context.Context, sessionID string, updates *ActiveSessionUpdates) error
	GetActiveSession(ctx context.Context, filters *ActiveSessionFilters) (*ActiveSessionData, error)
	ListActiveSessions(ctx context.Context, filters *ActiveSessionFilters) ([]*ActiveSessionData, error)
	DeleteActiveSession(ctx context.Context, sessionID string) error
	
	// Candle operations
	UpsertCandle(ctx context.Context, candle *CandleData) error
	GetCandles(ctx context.Context, filters *CandleFilters) ([]*CandleData, error)
	
	// Cleanup operations
	CleanupOldSessions(ctx context.Context) error
}

// CacheManager defines the interface for caching operations
type CacheManager interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, bool)
	Delete(key string) bool
	Clear() error
}

// Data transfer objects
type ActiveSessionData struct {
	ID        string
	SessionID string
	Exchange  string
	Symbol    string
	Type      string
	Timeframe *string
	Enabled   bool
}

type ActiveSessionUpdates struct {
	SessionID *string
	Enabled   *bool
}

type ActiveSessionFilters struct {
	Exchange  *string
	Symbol    *string
	Type      *string
	Timeframe *string
	Enabled   *bool
	SessionID *string
}

type CandleData struct {
	Exchange  string
	Symbol    string
	Timeframe string
	Timestamp int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type CandleFilters struct {
	Exchange  string
	Symbol    string
	Timeframe string
	Limit     int
}