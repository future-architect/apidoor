package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostAPIInfo(ctx context.Context, req *model.PostAPIInfoReq) error {
	if err := db.postAPIInfo(ctx, req); err != nil {
		log.Printf("db insert api info error: %v", err)
		return ServerError(err)
	}
	return nil
}
