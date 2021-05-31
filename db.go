package apidoor

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func GetAPIURL(key string, path string) error {
	if !rdb.SIsMember(ctx, key, path).Val() {
		return &MyError{Message: "unauthorized request"}
	}

	return nil
}
