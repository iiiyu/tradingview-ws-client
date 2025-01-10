package tvwsclient

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
)

// TVResponse represents the top-level response structure
type TVResponse struct {
	Method string        `json:"m"` // "qsd" for quote data
	Params []interface{} `json:"p"` // Array containing session ID and quote data
	Time   int64         `json:"t,omitempty"`
	TimeMS int64         `json:"t_ms,omitempty"`
}

// QuoteData represents the structure of quote data for a symbol
type QuoteData struct {
	Name   string     `json:"n"` // Symbol name (e.g., "BINANCE:BTCUSDT")
	Status string     `json:"s"` // Status ("ok")
	Values SymbolData `json:"v"` // Actual symbol data
}

// SymbolData represents the trading data for a symbol
type SymbolData struct {
	BaseCurrencyLogoID string   `json:"base-currency-logoid,omitempty"`
	Change             float64  `json:"ch,omitempty"`  // Price change
	ChangePercent      float64  `json:"chp,omitempty"` // Price change percentage
	CurrencyLogoID     string   `json:"currency-logoid,omitempty"`
	CurrencyCode       string   `json:"currency_code,omitempty"` // Currency (USD, USDT, etc.)
	CurrencyID         string   `json:"currency_id,omitempty"`
	BaseCurrencyID     string   `json:"base_currency_id,omitempty"`
	CurrentSession     string   `json:"current_session,omitempty"` // Session status
	Description        string   `json:"description,omitempty"`     // Full name
	Exchange           string   `json:"exchange,omitempty"`        // Exchange name
	Format             string   `json:"format,omitempty"`
	Fractional         bool     `json:"fractional,omitempty"`
	IsTradable         bool     `json:"is_tradable,omitempty"`
	Language           string   `json:"language,omitempty"`
	LocalDescription   string   `json:"local_description,omitempty"`
	ListedExchange     string   `json:"listed_exchange,omitempty"`
	LogoID             string   `json:"logoid,omitempty"`
	LastPrice          float64  `json:"lp,omitempty"`      // Last price
	LastPriceTime      int64    `json:"lp_time,omitempty"` // Timestamp
	MinMove            int      `json:"minmov,omitempty"`
	MinMove2           int      `json:"minmove2,omitempty"`
	OriginalName       string   `json:"original_name,omitempty"`
	PriceScale         int      `json:"pricescale,omitempty"` // Price scale factor
	ProName            string   `json:"pro_name,omitempty"`
	ShortName          string   `json:"short_name,omitempty"` // Symbol short name
	Type               string   `json:"type,omitempty"`       // Asset type (stock, spot, etc.)
	TypeSpecs          []string `json:"typespecs,omitempty"`
	UpdateMode         string   `json:"update_mode,omitempty"`
	Volume             float64  `json:"volume,omitempty"` // Trading volume
	VariableTickSize   string   `json:"variable_tick_size,omitempty"`
	ValueUnitID        string   `json:"value_unit_id,omitempty"`
	UnitID             string   `json:"unit_id,omitempty"`
	Measure            string   `json:"measure,omitempty"`
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
	SessionID string
	Symbol    string
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
