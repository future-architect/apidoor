package apiredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

type RedisDB struct {
	client *redis.Client
}

func NewRedisDB() *RedisDB {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	return &RedisDB{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}),
	}
}

func (rd RedisDB) PostAPIRouting(ctx context.Context, apikey, path, forwardURL string) error {
	err := rd.client.HSet(ctx, apikey, path, forwardURL).Err()
	return err
}
