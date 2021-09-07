package gateway

import (
	"context"
	"fmt"
	"log"
	"os"
)

var DBDriver DB

func init() {
	dbType := os.Getenv("API_DB_TYPE")
	var err error
	if DBDriver, err = createDBDriver(dbType); err != nil {
		log.Fatalf("api db set up error: %v", err)
	}
}

type DB interface {
	GetFields(ctx context.Context, key string) (Fields, error)
}

func createDBDriver(dbType string) (DB, error) {
	if dbType == "REDIS" {
		return NewRedisDB(), nil
	} else {
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}
