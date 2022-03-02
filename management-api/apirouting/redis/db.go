package redis

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/go-redis/redis/v8"
	"os"
)

type APIRouting struct {
	client *redis.Client
}

func New() *APIRouting {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	return &APIRouting{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
	}
}

func (ar APIRouting) PostRouting(ctx context.Context, apikey, path, forwardURL string) error {
	err := ar.client.HSet(ctx, apikey, path, forwardURL).Err()
	return err
}
func (ar APIRouting) CountRouting(ctx context.Context, apikey, path string) (int64, error) {
	//TODO: impl
	return 0, nil
}

func (ar APIRouting) PostAPIToken(ctx context.Context, req model.PostAPITokenReq) error {
	//TODO: impl
	return nil
}
