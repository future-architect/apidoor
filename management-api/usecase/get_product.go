package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func GetProducts(ctx context.Context) ([]model.Product, error) {
	list, err := db.getProducts(ctx)
	if err != nil {
		log.Printf("execute get product from db error: %v", err)
		return nil, err
	}
	return list, nil
}
