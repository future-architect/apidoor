package apidoor

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

func GetAPIURL(ctx context.Context, key, path string) (string, error) {
	if !rdb.HExists(ctx, key, path).Val() {
		return "", &MyError{Message: "unauthorized request"}
	}

	return rdb.HGet(ctx, key, path).Val()[1:], nil
}
