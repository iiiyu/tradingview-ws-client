package tvwsclient

// Option represents a client option
type Option func(*Client)

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
	BaseCurrencyLogoID string  `json:"base-currency-logoid"`
	Change             float64 `json:"ch"`  // Price change
	ChangePercent      float64 `json:"chp"` // Price change percentage
	CurrencyLogoID     string  `json:"currency-logoid"`
	CurrencyCode       string  `json:"currency_code"` // Currency (USD, USDT, etc.)
	CurrencyID         string  `json:"currency_id"`
	BaseCurrencyID     string  `json:"base_currency_id"`
	CurrentSession     string  `json:"current_session"` // Session status
	Description        string  `json:"description"`     // Full name
	Exchange           string  `json:"exchange"`        // Exchange name
	Format             string  `json:"format"`
	Fractional         bool    `json:"fractional"`
	IsTradable         bool    `json:"is_tradable"`
	Language           string  `json:"language"`
	LocalDescription   string  `json:"local_description"`
	ListedExchange     string  `json:"listed_exchange"`
	LogoID             string  `json:"logoid"`
	LastPrice          float64 `json:"lp"`      // Last price
	LastPriceTime      int64   `json:"lp_time"` // Timestamp
	MinMove            int     `json:"minmov"`
	MinMove2           int     `json:"minmove2"`
	OriginalName       string  `json:"original_name"`
	PriceScale         int     `json:"pricescale"` // Price scale factor
	ProName            string  `json:"pro_name"`
	ShortName          string  `json:"short_name"` // Symbol short name
	Type               string  `json:"type"`       // Asset type (stock, spot, etc.)
	TypeSpecs          string  `json:"typespecs"`
	UpdateMode         string  `json:"update_mode"`
	Volume             float64 `json:"volume"` // Trading volume
	VariableTickSize   bool    `json:"variable_tick_size"`
	ValueUnitID        string  `json:"value_unit_id"`
	UnitID             string  `json:"unit_id"`
	Measure            string  `json:"measure"`
}

// SeriesLoadingMessage represents the series_loading message
type SeriesLoadingMessage struct {
	ChartSessionID string
	SeriesID       string
	SeriesSet      string
}

// SymbolResolvedMessage represents the symbol_resolved message
type SymbolResolvedMessage struct {
	ChartSessionID string
	SeriesID       string
	SymbolInfo     SymbolInfo
}

// SymbolInfo contains detailed information about a symbol
type SymbolInfo struct {
	Source2           Source2Info `json:"source2"`
	CurrencyCode      string      `json:"currency_code"`
	SourceID          string      `json:"source_id"`
	SessionHolidays   string      `json:"session_holidays"`
	SubsessionID      string      `json:"subsession_id"`
	ProviderID        string      `json:"provider_id"`
	CurrencyID        string      `json:"currency_id"`
	Country           string      `json:"country"`
	ProPerm           string      `json:"pro_perm"`
	Measure           string      `json:"measure"`
	AllowedAdjustment string      `json:"allowed_adjustment"`
	ShortDescription  string      `json:"short_description"`
	VariableTickSize  string      `json:"variable_tick_size"`
	ISIN              string      `json:"isin"`
	Language          string      `json:"language"`
	LocalDescription  string      `json:"local_description"`
	Name              string      `json:"name"`
	FullName          string      `json:"full_name"`
	ProName           string      `json:"pro_name"`
	BaseName          []string    `json:"base_name"`
	Description       string      `json:"description"`
	Exchange          string      `json:"exchange"`
	PriceScale        int         `json:"pricescale"`
	PointValue        float64     `json:"pointvalue"`
	MinMove           int         `json:"minmov"`
	Session           string      `json:"session"`
	SessionDisplay    string      `json:"session_display"`
	Type              string      `json:"type"`
	HasIntraday       bool        `json:"has_intraday"`
	Fractional        bool        `json:"fractional"`
	ListedExchange    string      `json:"listed_exchange"`
	IsTradable        bool        `json:"is_tradable"`
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

// TimescaleUpdateMessage represents the timescale_update message
type TimescaleUpdateMessage struct {
	ChartSessionID string
	Data           TimescaleData
}

// TimescaleData contains the chart data
type TimescaleData struct {
	SeriesID string     `json:"sds_1"`
	Node     string     `json:"node"`
	Series   []ChartBar `json:"s"`
}

// ChartBar represents a single bar of chart data
type ChartBar struct {
	Index  int       `json:"i"`
	Values []float64 `json:"v"` // [timestamp, open, high, low, close, volume]
}

// QuoteCompletedMessage represents the quote_completed message
type QuoteCompletedMessage struct {
	SessionID string
	Symbol    string
}
