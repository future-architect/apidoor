package datasource

import (
	"context"
	"github.com/future-architect/apidoor/gateway/model"
)

type DataSource interface {
	GetFields(ctx context.Context, key string) (model.Fields, error)
}
