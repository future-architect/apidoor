package redis

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/gateway"
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

func (rd DataSource) GetFields(ctx context.Context, key string) (gateway.Fields, error) {
	var fields []gateway.Field

	for _, hk := range rd.client.HKeys(ctx, key).Val() {
		u := gateway.NewURITemplate(hk)
		v := gateway.NewURITemplate(rd.client.HGet(ctx, key, hk).Val())
		fields = append(fields, gateway.Field{
			Template: *u,
			Path:     *v,
			Num:      5,  // TODO
			Max:      10, // TODO
		})
	}

	if len(fields) == 0 {
		return nil, gateway.ErrUnauthorizedRequest
	}

	return fields, nil
}
