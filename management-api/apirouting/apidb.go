package apirouting

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/apirouting/dynamo"
	"github.com/future-architect/apidoor/managementapi/apirouting/redis"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
	"os"
)

var ApiDBDriver APIDB

func init() {
	dbType := os.Getenv("API_DB_TYPE")
	var err error
	if ApiDBDriver, err = createDBDriver(dbType); err != nil {
		log.Fatalf("api db set up error: %v", err)
	}
}

type APIDB interface {
	PostRouting(ctx context.Context, apiKey, path, forwardURL string) error
	PostAPIToken(ctx context.Context, req model.PostAPITokenReq) error
	CountRouting(ctx context.Context, apikey, path string) (int64, error)
}

func createDBDriver(dbType string) (APIDB, error) {
	switch dbType {
	case "REDIS":
		return redis.New(), nil
	case "DYNAMO":
		fallthrough
	case "":
		return dynamo.New(), nil
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}
