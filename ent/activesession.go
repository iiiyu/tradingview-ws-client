// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
)

// ActiveSession is the model entity for the ActiveSession schema.
type ActiveSession struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// SessionID holds the value of the "session_id" field.
	SessionID string `json:"session_id,omitempty"`
	// Exchange holds the value of the "exchange" field.
	Exchange string `json:"exchange,omitempty"`
	// Symbol holds the value of the "symbol" field.
	Symbol string `json:"symbol,omitempty"`
	// Timeframe holds the value of the "timeframe" field.
	Timeframe activesession.Timeframe `json:"timeframe,omitempty"`
	// Enabled holds the value of the "enabled" field.
	Enabled bool `json:"enabled,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ActiveSession) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case activesession.FieldEnabled:
			values[i] = new(sql.NullBool)
		case activesession.FieldID:
			values[i] = new(sql.NullInt64)
		case activesession.FieldSessionID, activesession.FieldExchange, activesession.FieldSymbol, activesession.FieldTimeframe:
			values[i] = new(sql.NullString)
		case activesession.FieldCreatedAt, activesession.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ActiveSession fields.
func (as *ActiveSession) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case activesession.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			as.ID = int(value.Int64)
		case activesession.FieldSessionID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field session_id", values[i])
			} else if value.Valid {
				as.SessionID = value.String
			}
		case activesession.FieldExchange:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field exchange", values[i])
			} else if value.Valid {
				as.Exchange = value.String
			}
		case activesession.FieldSymbol:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field symbol", values[i])
			} else if value.Valid {
				as.Symbol = value.String
			}
		case activesession.FieldTimeframe:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field timeframe", values[i])
			} else if value.Valid {
				as.Timeframe = activesession.Timeframe(value.String)
			}
		case activesession.FieldEnabled:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field enabled", values[i])
			} else if value.Valid {
				as.Enabled = value.Bool
			}
		case activesession.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				as.CreatedAt = value.Time
			}
		case activesession.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				as.UpdatedAt = value.Time
			}
		default:
			as.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ActiveSession.
// This includes values selected through modifiers, order, etc.
func (as *ActiveSession) Value(name string) (ent.Value, error) {
	return as.selectValues.Get(name)
}

// Update returns a builder for updating this ActiveSession.
// Note that you need to call ActiveSession.Unwrap() before calling this method if this ActiveSession
// was returned from a transaction, and the transaction was committed or rolled back.
func (as *ActiveSession) Update() *ActiveSessionUpdateOne {
	return NewActiveSessionClient(as.config).UpdateOne(as)
}

// Unwrap unwraps the ActiveSession entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (as *ActiveSession) Unwrap() *ActiveSession {
	_tx, ok := as.config.driver.(*txDriver)
	if !ok {
		panic("ent: ActiveSession is not a transactional entity")
	}
	as.config.driver = _tx.drv
	return as
}

// String implements the fmt.Stringer.
func (as *ActiveSession) String() string {
	var builder strings.Builder
	builder.WriteString("ActiveSession(")
	builder.WriteString(fmt.Sprintf("id=%v, ", as.ID))
	builder.WriteString("session_id=")
	builder.WriteString(as.SessionID)
	builder.WriteString(", ")
	builder.WriteString("exchange=")
	builder.WriteString(as.Exchange)
	builder.WriteString(", ")
	builder.WriteString("symbol=")
	builder.WriteString(as.Symbol)
	builder.WriteString(", ")
	builder.WriteString("timeframe=")
	builder.WriteString(fmt.Sprintf("%v", as.Timeframe))
	builder.WriteString(", ")
	builder.WriteString("enabled=")
	builder.WriteString(fmt.Sprintf("%v", as.Enabled))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(as.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(as.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ActiveSessions is a parsable slice of ActiveSession.
type ActiveSessions []*ActiveSession