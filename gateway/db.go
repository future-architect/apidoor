package gateway

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

type Field struct {
	Template URITemplate
	Path     URITemplate
	Max      int
}

type KeyData map[string][]Field

var APIData = make(KeyData)

func GetAPIURL(ctx context.Context, key, path string) (string, error) {
	fields, ok := APIData[key]
	if !ok {
		return "", &MyError{Message: "unauthorized request"}
	}

	u := NewURITemplate(path)
	for _, v := range fields {
		if _, ok := u.TemplateMatch(v.Template); ok {
			return v.Path.JoinPath(), nil
		}
	}

	return "", &MyError{Message: "unauthorized request"}
}

func init() {
	ctx := context.Background()
	for _, k := range rdb.Keys(ctx, "*").Val() {
		for _, hk := range rdb.HKeys(ctx, k).Val() {
			u := NewURITemplate(hk)
			v := NewURITemplate(rdb.HGet(ctx, k, hk).Val())
			n := 5
			/*
				n, err := strconv.Atoi(rdb.HGet(ctx, k, hk).Val())
				if err != nil {
					log.Fatal(err)
				}
			*/
			APIData[k] = append(APIData[k], Field{
				Template: *u,
				Path:     *v,
				Max:      n,
			})
		}
	}
}
