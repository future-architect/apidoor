package redis

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/gateway/datasource"
	"github.com/future-architect/apidoor/gateway/model"
	"github.com/go-redis/redis/v8"
	"os"
)

type DataSource struct {
	client *redis.Client
}

func New() *DataSource {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "6379"
	}

	return &DataSource{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: "",
			DB:       0,
		}),
	}
}

func (rd DataSource) GetFields(ctx context.Context, key string) (model.Fields, error) {
	var fields []model.Field

	for _, hk := range rd.client.HKeys(ctx, key).Val() {

		pathValue := rd.client.HGet(ctx, key, hk).Val()

		field, err := datasource.CreateField(ctx, key, hk, pathValue)
		if err != nil {
			return nil, fmt.Errorf("fetch field, key = %v, hk = %v, forwardURL = %v, error: %w",
				key, hk, pathValue, err)
		}

		fields = append(fields, field)
	}

	if len(fields) == 0 {
		return nil, model.ErrUnauthorizedRequest
	}

	return fields, nil
}
