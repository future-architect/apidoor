package gateway

import (
	"context"
	"os"
	"regexp"

	"github.com/go-redis/redis/v8"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

func GetAPIURL(ctx context.Context, key, path string) (string, error) {
	regexps := rdb.HKeys(ctx, key).Val()
	if len(regexps) == 0 {
		return "", &MyError{Message: "unauthorized request"}
	}

	for _, v := range regexps {
		re, err := regexp.Compile(v)
		if err != nil {
			return "", &MyError{Message: "unexpected field in redis"}
		}

		if re.Match([]byte(path)) {
			return rdb.HGet(ctx, key, v).Val()[1:], nil
		}
	}

	return "", &MyError{Message: "unauthorized request"}
}
