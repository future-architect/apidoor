package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func SearchAPIINfo(ctx context.Context, params *model.SearchAPIInfoParams) (*model.SearchAPIInfoResp, error) {
	respBody, err := db.searchAPIInfo(ctx, params)
	if err != nil {
		log.Printf("search api info db error: %v", err)
		return nil, ServerError{err}
	}
	return respBody, nil
}
