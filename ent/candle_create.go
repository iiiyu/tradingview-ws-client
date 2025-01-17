// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
)

// CandleCreate is the builder for creating a Candle entity.
type CandleCreate struct {
	config
	mutation *CandleMutation
	hooks    []Hook
}

// SetExchange sets the "exchange" field.
func (cc *CandleCreate) SetExchange(s string) *CandleCreate {
	cc.mutation.SetExchange(s)
	return cc
}

// SetSymbol sets the "symbol" field.
func (cc *CandleCreate) SetSymbol(s string) *CandleCreate {
	cc.mutation.SetSymbol(s)
	return cc
}

// SetTimeframe sets the "timeframe" field.
func (cc *CandleCreate) SetTimeframe(c candle.Timeframe) *CandleCreate {
	cc.mutation.SetTimeframe(c)
	return cc
}

// SetTimestamp sets the "timestamp" field.
func (cc *CandleCreate) SetTimestamp(i int64) *CandleCreate {
	cc.mutation.SetTimestamp(i)
	return cc
}

// SetOpen sets the "open" field.
func (cc *CandleCreate) SetOpen(f float64) *CandleCreate {
	cc.mutation.SetOpen(f)
	return cc
}

// SetHigh sets the "high" field.
func (cc *CandleCreate) SetHigh(f float64) *CandleCreate {
	cc.mutation.SetHigh(f)
	return cc
}

// SetLow sets the "low" field.
func (cc *CandleCreate) SetLow(f float64) *CandleCreate {
	cc.mutation.SetLow(f)
	return cc
}

// SetClose sets the "close" field.
func (cc *CandleCreate) SetClose(f float64) *CandleCreate {
	cc.mutation.SetClose(f)
	return cc
}

// SetVolume sets the "volume" field.
func (cc *CandleCreate) SetVolume(f float64) *CandleCreate {
	cc.mutation.SetVolume(f)
	return cc
}

// SetCreatedAt sets the "created_at" field.
func (cc *CandleCreate) SetCreatedAt(t time.Time) *CandleCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cc *CandleCreate) SetNillableCreatedAt(t *time.Time) *CandleCreate {
	if t != nil {
		cc.SetCreatedAt(*t)
	}
	return cc
}

// SetID sets the "id" field.
func (cc *CandleCreate) SetID(u uuid.UUID) *CandleCreate {
	cc.mutation.SetID(u)
	return cc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (cc *CandleCreate) SetNillableID(u *uuid.UUID) *CandleCreate {
	if u != nil {
		cc.SetID(*u)
	}
	return cc
}

// Mutation returns the CandleMutation object of the builder.
func (cc *CandleCreate) Mutation() *CandleMutation {
	return cc.mutation
}

// Save creates the Candle in the database.
func (cc *CandleCreate) Save(ctx context.Context) (*Candle, error) {
	cc.defaults()
	return withHooks(ctx, cc.sqlSave, cc.mutation, cc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CandleCreate) SaveX(ctx context.Context) *Candle {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CandleCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CandleCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *CandleCreate) defaults() {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		v := candle.DefaultCreatedAt()
		cc.mutation.SetCreatedAt(v)
	}
	if _, ok := cc.mutation.ID(); !ok {
		v := candle.DefaultID()
		cc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *CandleCreate) check() error {
	if _, ok := cc.mutation.Exchange(); !ok {
		return &ValidationError{Name: "exchange", err: errors.New(`ent: missing required field "Candle.exchange"`)}
	}
	if _, ok := cc.mutation.Symbol(); !ok {
		return &ValidationError{Name: "symbol", err: errors.New(`ent: missing required field "Candle.symbol"`)}
	}
	if _, ok := cc.mutation.Timeframe(); !ok {
		return &ValidationError{Name: "timeframe", err: errors.New(`ent: missing required field "Candle.timeframe"`)}
	}
	if v, ok := cc.mutation.Timeframe(); ok {
		if err := candle.TimeframeValidator(v); err != nil {
			return &ValidationError{Name: "timeframe", err: fmt.Errorf(`ent: validator failed for field "Candle.timeframe": %w`, err)}
		}
	}
	if _, ok := cc.mutation.Timestamp(); !ok {
		return &ValidationError{Name: "timestamp", err: errors.New(`ent: missing required field "Candle.timestamp"`)}
	}
	if _, ok := cc.mutation.Open(); !ok {
		return &ValidationError{Name: "open", err: errors.New(`ent: missing required field "Candle.open"`)}
	}
	if _, ok := cc.mutation.High(); !ok {
		return &ValidationError{Name: "high", err: errors.New(`ent: missing required field "Candle.high"`)}
	}
	if _, ok := cc.mutation.Low(); !ok {
		return &ValidationError{Name: "low", err: errors.New(`ent: missing required field "Candle.low"`)}
	}
	if _, ok := cc.mutation.Close(); !ok {
		return &ValidationError{Name: "close", err: errors.New(`ent: missing required field "Candle.close"`)}
	}
	if _, ok := cc.mutation.Volume(); !ok {
		return &ValidationError{Name: "volume", err: errors.New(`ent: missing required field "Candle.volume"`)}
	}
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Candle.created_at"`)}
	}
	return nil
}

func (cc *CandleCreate) sqlSave(ctx context.Context) (*Candle, error) {
	if err := cc.check(); err != nil {
		return nil, err
	}
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	cc.mutation.id = &_node.ID
	cc.mutation.done = true
	return _node, nil
}

func (cc *CandleCreate) createSpec() (*Candle, *sqlgraph.CreateSpec) {
	var (
		_node = &Candle{config: cc.config}
		_spec = sqlgraph.NewCreateSpec(candle.Table, sqlgraph.NewFieldSpec(candle.FieldID, field.TypeUUID))
	)
	if id, ok := cc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := cc.mutation.Exchange(); ok {
		_spec.SetField(candle.FieldExchange, field.TypeString, value)
		_node.Exchange = value
	}
	if value, ok := cc.mutation.Symbol(); ok {
		_spec.SetField(candle.FieldSymbol, field.TypeString, value)
		_node.Symbol = value
	}
	if value, ok := cc.mutation.Timeframe(); ok {
		_spec.SetField(candle.FieldTimeframe, field.TypeEnum, value)
		_node.Timeframe = value
	}
	if value, ok := cc.mutation.Timestamp(); ok {
		_spec.SetField(candle.FieldTimestamp, field.TypeInt64, value)
		_node.Timestamp = value
	}
	if value, ok := cc.mutation.Open(); ok {
		_spec.SetField(candle.FieldOpen, field.TypeFloat64, value)
		_node.Open = value
	}
	if value, ok := cc.mutation.High(); ok {
		_spec.SetField(candle.FieldHigh, field.TypeFloat64, value)
		_node.High = value
	}
	if value, ok := cc.mutation.Low(); ok {
		_spec.SetField(candle.FieldLow, field.TypeFloat64, value)
		_node.Low = value
	}
	if value, ok := cc.mutation.Close(); ok {
		_spec.SetField(candle.FieldClose, field.TypeFloat64, value)
		_node.Close = value
	}
	if value, ok := cc.mutation.Volume(); ok {
		_spec.SetField(candle.FieldVolume, field.TypeFloat64, value)
		_node.Volume = value
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.SetField(candle.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	return _node, _spec
}

// CandleCreateBulk is the builder for creating many Candle entities in bulk.
type CandleCreateBulk struct {
	config
	err      error
	builders []*CandleCreate
}

// Save creates the Candle entities in the database.
func (ccb *CandleCreateBulk) Save(ctx context.Context) ([]*Candle, error) {
	if ccb.err != nil {
		return nil, ccb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Candle, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CandleMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CandleCreateBulk) SaveX(ctx context.Context) []*Candle {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CandleCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CandleCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
