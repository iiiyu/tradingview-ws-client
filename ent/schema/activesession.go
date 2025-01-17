package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ActiveSession holds the schema definition for the ActiveSession entity.
type ActiveSession struct {
	ent.Schema
}

// Fields of the ActiveSession.
func (ActiveSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("session_id").Unique(),
		field.String("exchange"),
		field.String("symbol"),
		field.Enum("timeframe").Values("10S", "1", "5", "1D"),
		field.Bool("enabled").Default(false),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ActiveSession.
func (ActiveSession) Edges() []ent.Edge {
	return nil
}
