package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func DeleteAPIToken(ctx context.Context, req model.DeleteAPITokenReq) error {
	if err := apirouting.ApiDBDriver.DeleteAPIToken(ctx, req); err != nil {
		log.Printf("delete api token db error: %v", err)
		return ServerError{err}
	}
	return nil
}
