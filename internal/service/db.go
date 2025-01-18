package service

import (
	"context"

	"github.com/iiiyu/tradingview-ws-client/ent"
	"github.com/iiiyu/tradingview-ws-client/ent/activesession"
)

func CleanUpOldSessions(client *ent.Client) error {
	return client.ActiveSession.Update().
		Where(activesession.EnabledEQ(true)).
		SetEnabled(false).
		Exec(context.Background())
}
