package datasource

import (
	"context"
	"github.com/future-architect/apidoor/gateway/model"
)

type DataSource interface {
	GetFields(ctx context.Context, key string) (model.Fields, error)
	GetAccessTokens(ctx context.Context, apikey, templatePath string) (*model.AccessTokens, error)
}
