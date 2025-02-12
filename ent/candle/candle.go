// Code generated by ent, DO NOT EDIT.

package candle

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the candle type in the database.
	Label = "candle"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldExchange holds the string denoting the exchange field in the database.
	FieldExchange = "exchange"
	// FieldSymbol holds the string denoting the symbol field in the database.
	FieldSymbol = "symbol"
	// FieldTimeframe holds the string denoting the timeframe field in the database.
	FieldTimeframe = "timeframe"
	// FieldTimestamp holds the string denoting the timestamp field in the database.
	FieldTimestamp = "timestamp"
	// FieldOpen holds the string denoting the open field in the database.
	FieldOpen = "open"
	// FieldHigh holds the string denoting the high field in the database.
	FieldHigh = "high"
	// FieldLow holds the string denoting the low field in the database.
	FieldLow = "low"
	// FieldClose holds the string denoting the close field in the database.
	FieldClose = "close"
	// FieldVolume holds the string denoting the volume field in the database.
	FieldVolume = "volume"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// Table holds the table name of the candle in the database.
	Table = "candles"
)

// Columns holds all SQL columns for candle fields.
var Columns = []string{
	FieldID,
	FieldExchange,
	FieldSymbol,
	FieldTimeframe,
	FieldTimestamp,
	FieldOpen,
	FieldHigh,
	FieldLow,
	FieldClose,
	FieldVolume,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// Timeframe defines the type for the "timeframe" enum field.
type Timeframe string

// Timeframe values.
const (
	Timeframe10S Timeframe = "10S"
	Timeframe1   Timeframe = "1"
	Timeframe5   Timeframe = "5"
	Timeframe1D  Timeframe = "1D"
)

func (t Timeframe) String() string {
	return string(t)
}

// TimeframeValidator is a validator for the "timeframe" field enum values. It is called by the builders before save.
func TimeframeValidator(t Timeframe) error {
	switch t {
	case Timeframe10S, Timeframe1, Timeframe5, Timeframe1D:
		return nil
	default:
		return fmt.Errorf("candle: invalid enum value for timeframe field: %q", t)
	}
}

// OrderOption defines the ordering options for the Candle queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByExchange orders the results by the exchange field.
func ByExchange(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExchange, opts...).ToFunc()
}

// BySymbol orders the results by the symbol field.
func BySymbol(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSymbol, opts...).ToFunc()
}

// ByTimeframe orders the results by the timeframe field.
func ByTimeframe(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimeframe, opts...).ToFunc()
}

// ByTimestamp orders the results by the timestamp field.
func ByTimestamp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTimestamp, opts...).ToFunc()
}

// ByOpen orders the results by the open field.
func ByOpen(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldOpen, opts...).ToFunc()
}

// ByHigh orders the results by the high field.
func ByHigh(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHigh, opts...).ToFunc()
}

// ByLow orders the results by the low field.
func ByLow(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLow, opts...).ToFunc()
}

// ByClose orders the results by the close field.
func ByClose(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldClose, opts...).ToFunc()
}

// ByVolume orders the results by the volume field.
func ByVolume(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVolume, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}
