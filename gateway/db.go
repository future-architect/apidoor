package gateway

import (
	"context"
	"fmt"
	"log"
	"os"
)

var dbDriver DB

func init() {
	dbType := os.Getenv("DB_TYPE")
	var err error
	if dbDriver, err = createDBDriver(dbType); err != nil {
		log.Fatalf("db set up error: %v", err)
	}
}

type DB interface {
	GetFields(ctx context.Context, key string) (Fields, error)
}

func createDBDriver(dbType string) (DB, error) {
	switch dbType {
	case "REDIS":
		return NewRedisDB(), nil
	case "":
		// Using redis at default config
		return NewRedisDB(), nil
	case "DYNAMO":
		return NewDynamoDB(), nil
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}
