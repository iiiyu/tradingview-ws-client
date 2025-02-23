package tvwsclient

import (
	"encoding/json"
	"fmt"
)

// Option represents a client option
type Option func(*Client)

// WebSocket message methods
const (
	MethodQuoteData       = "qsd"
	MethodSeriesLoading   = "series_loading"
	MethodSymbolResolved  = "symbol_resolved"
	MethodTimescaleUpdate = "timescale_update"
	MethodSeriesCompleted = "series_completed"
	MethodDataUpdate      = "du"
	MethodQuoteCompleted  = "quote_completed"
)

// TVResponse represents the top-level response structure
type TVResponse struct {
	Method string        `json:"m"` // "qsd" for quote data
	Params []interface{} `json:"p"` // Array containing session ID and quote data
	Time   int64         `json:"t,omitempty"`
	TimeMS int64         `json:"t_ms,omitempty"`
}

type QuoteDataMessage struct {
	QuoteSessionID string    `json:"0"`
	Data           QuoteData `json:"1"`
}

// QuoteData represents the structure of quote data for a symbol
type QuoteData struct {
	Name   string     `json:"n"` // Symbol name with config (e.g., "={\"adjustment\":\"splits\",\"currency-id\":\"USD\",\"session\":\"regular\",\"symbol\":\"BATS:TSLA\"}")
	Status string     `json:"s"` // Status ("ok")
	Values SymbolData `json:"v"` // Actual symbol data
}

// SymbolData represents the trading data for a symbol
type SymbolData struct {
	BaseCurrencyLogoID  string   `json:"base-currency-logoid,omitempty"`
	Change              float64  `json:"ch,omitempty"`  // Price change
	ChangePercent       float64  `json:"chp,omitempty"` // Price change percentage
	CurrencyLogoID      string   `json:"currency-logoid,omitempty"`
	CurrencyCode        string   `json:"currency_code,omitempty"` // Currency (USD, USDT, etc.)
	CurrencyID          string   `json:"currency_id,omitempty"`
	BaseCurrencyID      string   `json:"base_currency_id,omitempty"`
	CurrentSession      string   `json:"current_session,omitempty"` // Session status
	Description         string   `json:"description,omitempty"`     // Full name
	Exchange            string   `json:"exchange,omitempty"`        // Exchange name
	Format              string   `json:"format,omitempty"`
	Fractional          bool     `json:"fractional,omitempty"`
	IsTradable          bool     `json:"is_tradable,omitempty"`
	Language            string   `json:"language,omitempty"`
	LocalDescription    string   `json:"local_description,omitempty"`
	ListedExchange      string   `json:"listed_exchange,omitempty"`
	LogoID              string   `json:"logoid,omitempty"`
	LastPrice           float64  `json:"lp,omitempty"`      // Last price
	LastPriceTime       int64    `json:"lp_time,omitempty"` // Timestamp
	MinMove             int      `json:"minmov,omitempty"`
	MinMove2            int      `json:"minmove2,omitempty"`
	OriginalName        string   `json:"original_name,omitempty"`
	PriceScale          int      `json:"pricescale,omitempty"` // Price scale factor
	ProName             string   `json:"pro_name,omitempty"`
	ShortName           string   `json:"short_name,omitempty"` // Symbol short name
	Type                string   `json:"type,omitempty"`       // Asset type (stock, spot, etc.)
	TypeSpecs           []string `json:"typespecs,omitempty"`
	UpdateMode          string   `json:"update_mode,omitempty"`
	Volume              float64  `json:"volume,omitempty"` // Trading volume
	VariableTickSize    string   `json:"variable_tick_size,omitempty"`
	ValueUnitID         string   `json:"value_unit_id,omitempty"`
	UnitID              string   `json:"unit_id,omitempty"`
	Measure             string   `json:"measure,omitempty"`
	RegularCloseTime    int64    `json:"regular_close_time,omitempty"`
	OpenTime            int64    `json:"open_time,omitempty"`
	RegularClose        float64  `json:"regular_close,omitempty"`
	OpenPrice           float64  `json:"open_price,omitempty"`
	HighPrice           float64  `json:"high_price,omitempty"`
	LowPrice            float64  `json:"low_price,omitempty"`
	AllTimeHigh         float64  `json:"all_time_high,omitempty"`
	AllTimeHighDay      int64    `json:"all_time_high_day,omitempty"`
	AllTimeLow          float64  `json:"all_time_low,omitempty"`
	AllTimeLowDay       int64    `json:"all_time_low_day,omitempty"`
	AverageVolume       float64  `json:"average_volume,omitempty"`
	Beta1Year           float64  `json:"beta_1_year,omitempty"`
	PriceEarningsTTM    float64  `json:"price_earnings_ttm,omitempty"`
	EarningsPerShareTTM float64  `json:"earnings_per_share_basic_ttm,omitempty"`
	EarningsPerShareFQ  float64  `json:"earnings_per_share_fq,omitempty"`
	MarketCapCalc       float64  `json:"market_cap_calc,omitempty"`
	TotalRevenue        float64  `json:"total_revenue,omitempty"`

	// Rates and broker information
	RatesMC     map[string]interface{} `json:"rates_mc,omitempty"`
	RatesTTM    map[string]interface{} `json:"rates_ttm,omitempty"`
	RatesFY     map[string]interface{} `json:"rates_fy,omitempty"`
	BrokerNames map[string]string      `json:"broker_names,omitempty"`

	// Options information
	OptionsInfo map[string]interface{} `json:"options-info,omitempty"`

	// Trading sessions
	Subsessions []SubsessionInfo `json:"subsessions,omitempty"`

	// Trade loaded - real time data
	TradeLoaded bool    `json:"trade_loaded,omitempty"`
	BidSize     float64 `json:"bid_size,omitempty"`
	Bid         float64 `json:"bid,omitempty"`
	AskSize     float64 `json:"ask_size,omitempty"`
	Ask         float64 `json:"ask,omitempty"`

	// RCH (Regular Change): The absolute price change during regular trading hours
	RCH float64 `json:"rch,omitempty"`
	// RCHP (Regular Change Percentage): The percentage change during regular trading hours
	RCHP float64 `json:"rchp,omitempty"`
	// RTC (Real-Time Close): The current/latest closing price in real-time
	RTC float64 `json:"rtc,omitempty"`
	// RTC_Time: The timestamp of the latest real-time close price
	RTC_Time int64 `json:"rtc_time,omitempty"`
}

// SeriesLoadingMessage represents the series_loading message
type SeriesLoadingMessage struct {
	ChartSessionID string       `json:"0,omitempty"` // First element in params array
	SeriesID       string       `json:"1,omitempty"` // Second element in params array
	SeriesSet      string       `json:"2,omitempty"` // Third element in params array
	SeriesNumber   string       `json:"3,omitempty"` // Fourth element in params array
	SeriesConfig   SeriesConfig `json:"4,omitempty"` // Fourth element in params array
}

// {"rt_update_period":0}
type SeriesConfig struct {
	RTUpdatePeriod int `json:"rt_update_period"` // 0
}

// SymbolResolvedMessage represents the symbol_resolved message
type SymbolResolvedMessage struct {
	ChartSessionID string     `json:"0"`
	SeriesID       string     `json:"1"`
	SymbolInfo     SymbolInfo `json:"2"`
}
type SymbolInfo struct {
	Source2             Source2Info      `json:"source2"`
	CurrencyCode        string           `json:"currency_code"`
	SourceID            string           `json:"source_id"`
	SessionHolidays     string           `json:"session_holidays"`
	SubsessionID        string           `json:"subsession_id"`
	ProviderID          string           `json:"provider_id"`
	CurrencyID          string           `json:"currency_id"`
	Country             string           `json:"country"`
	ProPerm             string           `json:"pro_perm"`
	Measure             string           `json:"measure"`
	AllowedAdjustment   string           `json:"allowed_adjustment"`
	ShortDescription    string           `json:"short_description"`
	VariableTickSize    string           `json:"variable_tick_size"`
	ISIN                string           `json:"isin"`
	Language            string           `json:"language"`
	LocalDescription    string           `json:"local_description"`
	Name                string           `json:"name"`
	FullName            string           `json:"full_name"`
	ProName             string           `json:"pro_name"`
	BaseName            []string         `json:"base_name"`
	Description         string           `json:"description"`
	Exchange            string           `json:"exchange"`
	PriceScale          int              `json:"pricescale"`
	PointValue          float64          `json:"pointvalue"`
	MinMove             int              `json:"minmov"`
	Session             string           `json:"session"`
	SessionDisplay      string           `json:"session_display"`
	Subsessions         []SubsessionInfo `json:"subsessions"`
	Type                string           `json:"type"`
	TypeSpecs           []string         `json:"typespecs"`
	HasIntraday         bool             `json:"has_intraday"`
	Fractional          bool             `json:"fractional"`
	ListedExchange      string           `json:"listed_exchange"`
	Legs                []string         `json:"legs"`
	IsTradable          bool             `json:"is_tradable"`
	MinMove2            int              `json:"minmove2"`
	Timezone            string           `json:"timezone"`
	Aliases             []string         `json:"aliases"`
	Alternatives        []string         `json:"alternatives"`
	IsReplayable        bool             `json:"is_replayable"`
	HasAdjustment       bool             `json:"has_adjustment"`
	HasExtendedHours    bool             `json:"has_extended_hours"`
	BarSource           string           `json:"bar_source"`
	BarTransform        string           `json:"bar_transform"`
	BarFillgaps         bool             `json:"bar_fillgaps"`
	VisiblePlotsSet     string           `json:"visible_plots_set"`
	IsTickbarsAvailable bool             `json:"is-tickbars-available"`
	FIGI                FIGIInfo         `json:"figi"`
}

// Add these new structs
type SubsessionInfo struct {
	Description       string `json:"description"`
	ID                string `json:"id"`
	Private           bool   `json:"private"`
	Session           string `json:"session"`
	SessionCorrection string `json:"session-correction,omitempty"`
	SessionDisplay    string `json:"session-display"`
}

type FIGIInfo struct {
	CountryComposite string `json:"country-composite"`
	ExchangeLevel    string `json:"exchange-level"`
}

// Source2Info contains information about the data source
type Source2Info struct {
	Country      string `json:"country"`
	Description  string `json:"description"`
	ExchangeType string `json:"exchange-type"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
}

type TimescaleUpdateMessage struct {
	ChartSessionID string              `json:"0"`
	Data           TimescaleUpdateData `json:"1"`
}

type TimescaleUpdateData struct {
	SDS1 struct {
		Node string `json:"node"`
		S    []struct {
			I int       `json:"i"`
			V []float64 `json:"v"` // [timestamp, open, high, low, close, volume]
		} `json:"s"`
		NS struct {
			D       string      `json:"d"`
			Indexes interface{} `json:"indexes"` // Can be string "nochange" or array []
		} `json:"ns"`
		T   string `json:"t"`
		Lbs struct {
			BarCloseTime int64 `json:"bar_close_time"`
		} `json:"lbs"`
	} `json:"sds_1"`
	Index     int     `json:"index"`
	Zoffset   int     `json:"zoffset"`
	Changes   []int64 `json:"changes"`
	Marks     [][]int `json:"marks"`
	IndexDiff []any   `json:"index_diff"`
	T         int64   `json:"t"`
	TMs       int64   `json:"t_ms"`
}

// QuoteCompletedMessage represents the quote_completed message
type QuoteCompletedMessage struct {
	SessionID       string
	ReceivedMessage string
}

// SeriesCompletedMessage represents the series_completed message
type SeriesCompletedMessage struct {
	ChartSessionID string       // "cs_Djf7086hIqtS"
	SeriesID       string       // "sds_1"
	Status         string       // "streaming"
	SeriesSet      string       // "s1"
	Config         SeriesConfig // Contains configuration parameters
	Time           int64        `json:"t,omitempty"`    // 1736302609
	TimeMS         int64        `json:"t_ms,omitempty"` // 1736302609050
}

// DuMessage represents the data update message structure
type DuMessage struct {
	ChartSessionID string `json:"0"` // "cs_lZqOBD1Jtvjb"
	Data           DuData `json:"1"` // The nested data object
}

type DuData struct {
	SDS1 struct {
		LBS struct {
			BarCloseTime int64 `json:"bar_close_time"`
		} `json:"lbs"`
		NS struct {
			D       string      `json:"d"`
			Indexes interface{} `json:"indexes"` // Can be "nochange" or other values
		} `json:"ns"`
		S []DuSeriesData `json:"s"`
		T string         `json:"t"` // Series set (e.g., "s1")
	} `json:"sds_1"`
}

type DuSeriesData struct {
	I int       `json:"i"` // Index
	V []float64 `json:"v"` // [timestamp, open, high, low, close, volume]
}

func NewQuoteCompletedMessage(params []interface{}) (*QuoteCompletedMessage, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("insufficient parameters")
	}

	sessionID, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid session ID type")
	}

	receivedMessage, ok := params[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid symbol type")
	}

	return &QuoteCompletedMessage{
		SessionID:       sessionID,
		ReceivedMessage: receivedMessage,
	}, nil
}

func NewQuoteDataMessage(params []interface{}) (*QuoteDataMessage, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("insufficient parameters")
	}

	sessionID, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid session ID type")
	}

	// Convert the interface{} to QuoteData through JSON marshaling
	paramJSON, err := json.Marshal(params[1])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal param: %w", err)
	}

	var quoteData QuoteData
	if err := json.Unmarshal(paramJSON, &quoteData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quote data: %w", err)
	}

	return &QuoteDataMessage{
		QuoteSessionID: sessionID,
		Data:           quoteData,
	}, nil
}

func NewSeriesLoadingMessage(params []interface{}) (*SeriesLoadingMessage, error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("insufficient parameters: expected at least 3, got %d", len(params))
	}

	// Convert required fields
	chartSessionID, ok1 := params[0].(string)
	seriesID, ok2 := params[1].(string)
	seriesSet, ok3 := params[2].(string)

	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid parameter types for required fields")
	}

	message := &SeriesLoadingMessage{
		ChartSessionID: chartSessionID,
		SeriesID:       seriesID,
		SeriesSet:      seriesSet,
	}

	// Handle optional SeriesNumber
	if len(params) >= 4 {
		if seriesNumber, ok := params[3].(string); ok {
			message.SeriesNumber = seriesNumber
		}
	}

	// Handle optional SeriesConfig
	if len(params) >= 5 {
		if configData, ok := params[4].(map[string]interface{}); ok {
			if period, ok := configData["rt_update_period"].(float64); ok {
				message.SeriesConfig = SeriesConfig{
					RTUpdatePeriod: int(period),
				}
			}
		}
	}

	return message, nil
}

func NewSymbolResolvedMessage(params []interface{}) (*SymbolResolvedMessage, error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("insufficient parameters: expected at least 3, got %d", len(params))
	}

	// Extract chart session ID and series ID
	chartSessionID, ok1 := params[0].(string)
	seriesID, ok2 := params[1].(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid parameter types for session ID or series ID")
	}

	// Convert the symbol info interface{} to SymbolInfo through JSON marshaling
	paramJSON, err := json.Marshal(params[2])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symbol info param: %w", err)
	}

	var symbolInfo SymbolInfo
	if err := json.Unmarshal(paramJSON, &symbolInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal symbol info: %w", err)
	}

	return &SymbolResolvedMessage{
		ChartSessionID: chartSessionID,
		SeriesID:       seriesID,
		SymbolInfo:     symbolInfo,
	}, nil
}

func NewTimescaleUpdateMessage(params []interface{}) (*TimescaleUpdateMessage, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("insufficient parameters: expected at least 2, got %d", len(params))
	}

	// Extract chart session ID
	chartSessionID, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid parameter type for chart session ID")
	}

	// Convert the data interface{} to TimescaleUpdateData through JSON marshaling
	paramJSON, err := json.Marshal(params[1])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal timescale update param: %w", err)
	}

	var timescaleUpdateData TimescaleUpdateData
	if err := json.Unmarshal(paramJSON, &timescaleUpdateData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal timescale update data: %w", err)
	}

	return &TimescaleUpdateMessage{
		ChartSessionID: chartSessionID,
		Data:           timescaleUpdateData,
	}, nil
}

func NewSeriesCompletedMessage(params []interface{}) (*SeriesCompletedMessage, error) {
	if len(params) < 5 {
		return nil, fmt.Errorf("insufficient parameters: expected at least 5, got %d", len(params))
	}

	// Extract required fields
	chartSessionID, ok1 := params[0].(string)
	seriesID, ok2 := params[1].(string)
	status, ok3 := params[2].(string)
	seriesSet, ok4 := params[3].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, fmt.Errorf("invalid parameter types for required fields")
	}

	message := &SeriesCompletedMessage{
		ChartSessionID: chartSessionID,
		SeriesID:       seriesID,
		Status:         status,
		SeriesSet:      seriesSet,
	}

	// Handle the SeriesConfig
	if configData, ok := params[4].(map[string]interface{}); ok {
		if period, ok := configData["rt_update_period"].(float64); ok {
			message.Config = SeriesConfig{
				RTUpdatePeriod: int(period),
			}
		}
	}

	return message, nil
}

func NewDuMessage(params []interface{}) (*DuMessage, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("insufficient parameters: expected at least 2, got %d", len(params))
	}

	// Extract chart session ID
	chartSessionID, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid parameter type for chart session ID")
	}

	// Convert the data interface{} to DuData through JSON marshaling
	paramJSON, err := json.Marshal(params[1])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal du data param: %w", err)
	}

	var duData DuData
	if err := json.Unmarshal(paramJSON, &duData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal du data: %w", err)
	}

	return &DuMessage{
		ChartSessionID: chartSessionID,
		Data:           duData,
	}, nil
}
