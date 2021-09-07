package gateway

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

func (rd RedisDB) GetFields(ctx context.Context, key string) (Fields, error) {
	var fields []Field

	for _, hk := range rd.client.HKeys(ctx, key).Val() {
		u := NewURITemplate(hk)
		v := NewURITemplate(rd.client.HGet(ctx, key, hk).Val())
		fields = append(fields, Field{
			Template: *u,
			Path:     *v,
			Num:      5,  // TODO
			Max:      10, // TODO
		})
	}

	if len(fields) == 0 {
		return nil, &MyError{Message: "unauthorized request"}
	}

	return fields, nil
}
