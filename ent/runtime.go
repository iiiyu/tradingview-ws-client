// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
	"github.com/iiiyu/tradingview-ws-client/ent/candle"
	"github.com/iiiyu/tradingview-ws-client/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	activesessionFields := schema.ActiveSession{}.Fields()
	_ = activesessionFields
	// activesessionDescEnabled is the schema descriptor for enabled field.
	activesessionDescEnabled := activesessionFields[4].Descriptor()
	// activesession.DefaultEnabled holds the default value on creation for the enabled field.
	activesession.DefaultEnabled = activesessionDescEnabled.Default.(bool)
	// activesessionDescCreatedAt is the schema descriptor for created_at field.
	activesessionDescCreatedAt := activesessionFields[5].Descriptor()
	// activesession.DefaultCreatedAt holds the default value on creation for the created_at field.
	activesession.DefaultCreatedAt = activesessionDescCreatedAt.Default.(func() time.Time)
	// activesessionDescUpdatedAt is the schema descriptor for updated_at field.
	activesessionDescUpdatedAt := activesessionFields[6].Descriptor()
	// activesession.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	activesession.DefaultUpdatedAt = activesessionDescUpdatedAt.Default.(func() time.Time)
	// activesession.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	activesession.UpdateDefaultUpdatedAt = activesessionDescUpdatedAt.UpdateDefault.(func() time.Time)
	candleFields := schema.Candle{}.Fields()
	_ = candleFields
	// candleDescCreatedAt is the schema descriptor for created_at field.
	candleDescCreatedAt := candleFields[9].Descriptor()
	// candle.DefaultCreatedAt holds the default value on creation for the created_at field.
	candle.DefaultCreatedAt = candleDescCreatedAt.Default.(func() time.Time)
}
