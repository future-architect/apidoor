package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostRouting(ctx context.Context, req model.PostAPIRoutingReq) error {
	if err := apirouting.ApiDBDriver.PostRouting(ctx, req.ApiKey, req.Path, req.ForwardURL); err != nil {
		log.Printf("post api routing db error: %v", err)
		return ServerError{err}
	}
	return nil
}
