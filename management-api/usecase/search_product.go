package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func SearchProduct(ctx context.Context, params *model.SearchProductParams) (*model.SearchProductResp, error) {
	respBody, err := db.searchProduct(ctx, params)
	if err != nil {
		log.Printf("search product db error: %v", err)
		return nil, ServerError{err}
	}
	return respBody, nil
}
