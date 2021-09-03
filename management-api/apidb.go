package managementapi

import (
	"context"
	"fmt"
	"log"
	"managementapi/apiredis"
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
	PostAPIRouting(ctx context.Context, apiKey, path, forwardURL string) error
}

func createDBDriver(dbType string) (APIDB, error) {
	if dbType == "REDIS" {
		return apiredis.NewRedisDB(), nil
	} else {
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}
