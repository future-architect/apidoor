package managementapi

import (
	"context"
	"log"
	"os"
)

var db sqlDB

func init() {
	dbDriver := os.Getenv("DATABASE_DRIVER")
	var err error
	switch dbDriver {
	case "postgres":
		if db, err = NewPostgresDB(); err != nil {
			log.Fatalf("setup postgreSQL failed: %v", err)
		}
	default:
		log.Fatalf("DATABASE_DRIVER is empty or not supported: %s", dbDriver)
	}
}

type sqlDB interface {
	getProducts(ctx context.Context) ([]Product, error)
	postProducts(ctx context.Context, product *PostProductReq) error
	searchProducts(ctx context.Context, params *SearchProductsParams) (*SearchProductsResp, error)
}
