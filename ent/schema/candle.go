package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Candle holds the schema definition for the Candle entity.
type Candle struct {
	ent.Schema
}

// Fields of the Candle.
func (Candle) Fields() []ent.Field {
	return []ent.Field{
		field.String("exchange"),
		field.String("symbol"),
		field.Enum("timeframe").Values("10S", "1", "5", "1D"),
		field.Int64("timestamp"),
		field.Float("open").GoType(float64(0)),
		field.Float("high").GoType(float64(0)),
		field.Float("low").GoType(float64(0)),
		field.Float("close").GoType(float64(0)),
		field.Float("volume").GoType(float64(0)),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Candle.
func (Candle) Edges() []ent.Edge {
	return nil
}

// Indexes of the Candle.
func (Candle) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("exchange", "symbol"),
		index.Fields("timestamp"),
		index.Fields("exchange", "symbol", "timeframe", "timestamp").Unique(),
	}
}
