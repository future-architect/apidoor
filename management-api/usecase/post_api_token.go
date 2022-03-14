package usecase

import (
	"context"
	"errors"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostAPIToken(ctx context.Context, req model.PostAPITokenReq) error {
	// check whether api routing exists
	cnt, err := apirouting.ApiDBDriver.CountRouting(ctx, req.APIKey, req.Path)
	if err != nil {
		log.Printf("count api routings db error: %v", err)
		return ServerError{err}
	}
	if cnt == 0 {
		log.Println("api_key or path is wrong")
		return ClientError{errors.New("api_key or path is wrong")}
	}

	if err := apirouting.ApiDBDriver.PostAPIToken(ctx, req); err != nil {
		log.Printf("insert api token db error: %v", err)
		return ServerError{err}
	}
	return nil
}
