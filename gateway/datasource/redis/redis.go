package redis

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/gateway/model"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
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

		var schema string
		if strings.HasPrefix(pathValue, "http://") {
			schema = "http"
			pathValue = strings.Replace(pathValue, "http://", "", 1)
			fmt.Println("Redis After:", pathValue)
		} else if strings.HasPrefix(pathValue, "https://") {
			schema = "https"
			pathValue = strings.Replace(pathValue, "https://", "", 1)
		} else {
			// スキーマが存在しない(tcpなどのスキーマは非対応)
			schema = "http"
		}

		fmt.Printf("%+v\n", model.NewURITemplate(pathValue))

		fields = append(fields, model.Field{
			Template:      model.NewURITemplate(hk),
			ForwardSchema: schema,
			Path:          model.NewURITemplate(pathValue),
			Num:           5,  // TODO
			Max:           10, // TODO
		})
	}

	if len(fields) == 0 {
		return nil, model.ErrUnauthorizedRequest
	}

	return fields, nil
}
