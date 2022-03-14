package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostUser(ctx context.Context, req model.PostUserReq) error {
	if err := db.postUser(ctx, &req); err != nil {
		log.Printf("db insert user error: %v", err)
		return ServerError{err}
	}
	return nil
}
