// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/ent/predicate"
)

// CandleDelete is the builder for deleting a Candle entity.
type CandleDelete struct {
	config
	hooks    []Hook
	mutation *CandleMutation
}

// Where appends a list predicates to the CandleDelete builder.
func (cd *CandleDelete) Where(ps ...predicate.Candle) *CandleDelete {
	cd.mutation.Where(ps...)
	return cd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cd *CandleDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cd.sqlExec, cd.mutation, cd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cd *CandleDelete) ExecX(ctx context.Context) int {
	n, err := cd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cd *CandleDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(candle.Table, sqlgraph.NewFieldSpec(candle.FieldID, field.TypeInt))
	if ps := cd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cd.mutation.done = true
	return affected, err
}

// CandleDeleteOne is the builder for deleting a single Candle entity.
type CandleDeleteOne struct {
	cd *CandleDelete
}

// Where appends a list predicates to the CandleDelete builder.
func (cdo *CandleDeleteOne) Where(ps ...predicate.Candle) *CandleDeleteOne {
	cdo.cd.mutation.Where(ps...)
	return cdo
}

// Exec executes the deletion query.
func (cdo *CandleDeleteOne) Exec(ctx context.Context) error {
	n, err := cdo.cd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{candle.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cdo *CandleDeleteOne) ExecX(ctx context.Context) {
	if err := cdo.Exec(ctx); err != nil {
		panic(err)
	}
}