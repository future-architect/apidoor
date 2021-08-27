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

func GetFields(ctx context.Context, key string) (Fields, error) {
	var fields []Field

	for _, hk := range rdb.HKeys(ctx, key).Val() {
		u := NewURITemplate(hk)
		v := NewURITemplate(rdb.HGet(ctx, key, hk).Val())
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


//func GetURL(ctx context.Context, key, path string) (string, error) {
//	var apiMapping []Field
//
//	for _, hk := range rdb.HKeys(ctx, key).Val() {
//		u := NewURITemplate(hk)
//		v := NewURITemplate(rdb.HGet(ctx, key, hk).Val())
//		apiMapping= append(apiMapping, Field{
//			Template: *u,
//			Path:     *v,
//			Num:      5,  // TODO
//			Max:      10, // TODO
//		})
//	}
//
//	if len(apiMapping) == 0 {
//		return "", &MyError{Message: "unauthorized request"}
//	}
//
//	u := NewURITemplate(path)
//	for _, v := range apiMapping {
//		if _, ok := u.TemplateMatch(v.Template); ok {
//			return v.Path.JoinPath(), nil
//		}
//	}
//
//	return "", &MyError{Message: "unauthorized request"}
//}
