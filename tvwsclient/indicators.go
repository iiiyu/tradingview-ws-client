package tvwsclient

import (
	"context"
	"fmt"
	"log/slog"
)

// IndicatorType represents different types of technical indicators
type IndicatorType string

const (
	IndicatorRSI              IndicatorType = "RSI@tv-basicstudies-1"
	IndicatorMACD             IndicatorType = "MACD@tv-basicstudies-1"
	IndicatorBollingerBands   IndicatorType = "BB@tv-basicstudies-1"
	IndicatorMovingAverage    IndicatorType = "MASimple@tv-basicstudies-1"
	IndicatorEMA              IndicatorType = "MAExp@tv-basicstudies-1"
	IndicatorVolume           IndicatorType = "Volume@tv-basicstudies-1"
	IndicatorStochastic       IndicatorType = "Stochastic@tv-basicstudies-1"
	IndicatorWilliamsR        IndicatorType = "WilliamsR@tv-basicstudies-1"
	IndicatorCCI              IndicatorType = "CCI@tv-basicstudies-1"
	IndicatorMomentum         IndicatorType = "Mom@tv-basicstudies-1"
)

// IndicatorConfig represents configuration for a technical indicator
type IndicatorConfig struct {
	Type       IndicatorType          `json:"type"`
	Version    string                 `json:"version"`
	Parameters map[string]interface{} `json:"parameters"`
	Enabled    bool                   `json:"enabled"`
}

// StudySession represents a study session with multiple indicators
type StudySession struct {
	SessionID  string                     `json:"session_id"`
	Symbol     string                     `json:"symbol"`
	Interval   string                     `json:"interval"`
	Indicators map[string]IndicatorConfig `json:"indicators"`
	Enabled    bool                       `json:"enabled"`
}

// StudyDataMessage represents incoming study/indicator data
type StudyDataMessage struct {
	StudySessionID string                 `json:"study_session_id"`
	StudyID        string                 `json:"study_id"`
	SeriesID       string                 `json:"series_id"`
	Timestamp      int64                  `json:"timestamp"`
	Values         map[string]interface{} `json:"values"`
	Type           string                 `json:"type"`
}

// IndicatorValue represents a calculated indicator value
type IndicatorValue struct {
	IndicatorType IndicatorType `json:"indicator_type"`
	Timestamp     int64         `json:"timestamp"`
	Values        map[string]float64 `json:"values"`
}

// StudyManager handles technical indicator sessions
type StudyManager struct {
	client   TradingViewClient
	sessions map[string]*StudySession
	logger   *slog.Logger
}

// NewStudyManager creates a new study manager
func NewStudyManager(client TradingViewClient, logger *slog.Logger) *StudyManager {
	return &StudyManager{
		client:   client,
		sessions: make(map[string]*StudySession),
		logger:   logger,
	}
}

// CreateStudySession creates a new study session for indicators
func (sm *StudyManager) CreateStudySession(symbol, interval string) (*StudySession, error) {
	sessionID := GenerateSession("study_")
	
	session := &StudySession{
		SessionID:  sessionID,
		Symbol:     symbol,
		Interval:   interval,
		Indicators: make(map[string]IndicatorConfig),
		Enabled:    true,
	}
	
	sm.sessions[sessionID] = session
	
	// Send create session message
	if err := sm.sendCreateStudySession(sessionID, symbol, interval); err != nil {
		delete(sm.sessions, sessionID)
		return nil, WrapMessageError("create_study_session", err)
	}
	
	sm.logger.Info("created study session", 
		"session_id", sessionID,
		"symbol", symbol,
		"interval", interval)
	
	return session, nil
}

// AddIndicator adds a technical indicator to a study session
func (sm *StudyManager) AddIndicator(sessionID string, indicatorID string, config IndicatorConfig) error {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return NewTradingViewError("add_indicator", ErrCodeSession, "study session not found", ErrSessionNotFound)
	}
	
	session.Indicators[indicatorID] = config
	
	// Send add indicator message
	if err := sm.sendAddIndicator(sessionID, indicatorID, config); err != nil {
		delete(session.Indicators, indicatorID)
		return WrapMessageError("add_indicator", err)
	}
	
	sm.logger.Info("added indicator to study session",
		"session_id", sessionID,
		"indicator_id", indicatorID,
		"indicator_type", config.Type)
	
	return nil
}

// RemoveIndicator removes an indicator from a study session
func (sm *StudyManager) RemoveIndicator(sessionID, indicatorID string) error {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return NewTradingViewError("remove_indicator", ErrCodeSession, "study session not found", ErrSessionNotFound)
	}
	
	delete(session.Indicators, indicatorID)
	
	// Send remove indicator message
	if err := sm.sendRemoveIndicator(sessionID, indicatorID); err != nil {
		return WrapMessageError("remove_indicator", err)
	}
	
	sm.logger.Info("removed indicator from study session",
		"session_id", sessionID,
		"indicator_id", indicatorID)
	
	return nil
}

// GetSession returns a study session by ID
func (sm *StudyManager) GetSession(sessionID string) (*StudySession, bool) {
	session, exists := sm.sessions[sessionID]
	return session, exists
}

// ListSessions returns all active study sessions
func (sm *StudyManager) ListSessions() []*StudySession {
	sessions := make([]*StudySession, 0, len(sm.sessions))
	for _, session := range sm.sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// DeleteSession removes a study session
func (sm *StudyManager) DeleteSession(sessionID string) error {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return NewTradingViewError("delete_session", ErrCodeSession, "study session not found", ErrSessionNotFound)
	}
	
	// Remove all indicators first
	for indicatorID := range session.Indicators {
		if err := sm.RemoveIndicator(sessionID, indicatorID); err != nil {
			sm.logger.Error("failed to remove indicator during session deletion",
				"session_id", sessionID,
				"indicator_id", indicatorID,
				"error", err)
		}
	}
	
	// Delete the session
	delete(sm.sessions, sessionID)
	
	sm.logger.Info("deleted study session", "session_id", sessionID)
	return nil
}

// ProcessStudyData processes incoming study data messages
func (sm *StudyManager) ProcessStudyData(ctx context.Context, msg *StudyDataMessage) error {
	sm.logger.Debug("processing study data",
		"study_session_id", msg.StudySessionID,
		"study_id", msg.StudyID,
		"timestamp", msg.Timestamp)
	
	// Find the corresponding session
	session, exists := sm.sessions[msg.StudySessionID]
	if !exists {
		sm.logger.Warn("received study data for unknown session", "session_id", msg.StudySessionID)
		return nil
	}
	
	// Process indicator values
	return sm.processIndicatorValues(session, msg)
}

// Default indicator configurations
func GetDefaultIndicatorConfig(indicatorType IndicatorType) IndicatorConfig {
	configs := map[IndicatorType]IndicatorConfig{
		IndicatorRSI: {
			Type:    IndicatorRSI,
			Version: "1",
			Parameters: map[string]interface{}{
				"length": 14,
			},
			Enabled: true,
		},
		IndicatorMACD: {
			Type:    IndicatorMACD,
			Version: "1",
			Parameters: map[string]interface{}{
				"fast_length":   12,
				"slow_length":   26,
				"signal_length": 9,
			},
			Enabled: true,
		},
		IndicatorBollingerBands: {
			Type:    IndicatorBollingerBands,
			Version: "1",
			Parameters: map[string]interface{}{
				"length": 20,
				"mult":   2.0,
			},
			Enabled: true,
		},
		IndicatorMovingAverage: {
			Type:    IndicatorMovingAverage,
			Version: "1",
			Parameters: map[string]interface{}{
				"length": 20,
			},
			Enabled: true,
		},
		IndicatorEMA: {
			Type:    IndicatorEMA,
			Version: "1",
			Parameters: map[string]interface{}{
				"length": 20,
			},
			Enabled: true,
		},
	}
	
	if config, exists := configs[indicatorType]; exists {
		return config
	}
	
	// Default fallback
	return IndicatorConfig{
		Type:       indicatorType,
		Version:    "1",
		Parameters: make(map[string]interface{}),
		Enabled:    true,
	}
}

// Private methods for WebSocket communication
func (sm *StudyManager) sendCreateStudySession(sessionID, symbol, interval string) error {
	// This would send the appropriate WebSocket message to create a study session
	// Implementation depends on TradingView's WebSocket protocol for studies
	sm.logger.Debug("sending create study session message",
		"session_id", sessionID,
		"symbol", symbol,
		"interval", interval)
	return nil
}

func (sm *StudyManager) sendAddIndicator(sessionID, indicatorID string, config IndicatorConfig) error {
	// This would send the appropriate WebSocket message to add an indicator
	sm.logger.Debug("sending add indicator message",
		"session_id", sessionID,
		"indicator_id", indicatorID,
		"indicator_type", config.Type)
	return nil
}

func (sm *StudyManager) sendRemoveIndicator(sessionID, indicatorID string) error {
	// This would send the appropriate WebSocket message to remove an indicator
	sm.logger.Debug("sending remove indicator message",
		"session_id", sessionID,
		"indicator_id", indicatorID)
	return nil
}

func (sm *StudyManager) processIndicatorValues(session *StudySession, msg *StudyDataMessage) error {
	// Process and store indicator values
	// This could write to database, cache, or send to API clients
	sm.logger.Debug("processing indicator values",
		"session_id", session.SessionID,
		"symbol", session.Symbol,
		"values", msg.Values)
	return nil
}

// WebSocket message constants for studies
const (
	MethodStudyLoading     = "study_loading"
	MethodStudyCompleted   = "study_completed"
	MethodStudyData        = "study_data"
	MethodStudyError       = "study_error"
	MethodSeriesStudyData  = "series_study_data"
)

// NewStudyDataMessage creates a StudyDataMessage from WebSocket parameters
func NewStudyDataMessage(params []interface{}) (*StudyDataMessage, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("insufficient parameters for study data message")
	}
	
	studySessionID, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid study session ID type")
	}
	
	// Parse study data from params[1]
	// This would need to be implemented based on TradingView's actual message format
	
	return &StudyDataMessage{
		StudySessionID: studySessionID,
		// ... other fields would be parsed from params
	}, nil
}