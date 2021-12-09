package datasource

import (
	"context"
	"github.com/future-architect/apidoor/gateway"
)

type DataSource interface {
	GetFields(ctx context.Context, key string) (gateway.Fields, error)
}
