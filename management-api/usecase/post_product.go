package usecase

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostProduct(ctx context.Context, req model.PostProductReq) error {
	if err := db.postProduct(ctx, &req); err != nil {
		log.Printf("insert product to db failed: %v", err)
		if constraintErr, ok := err.(*dbConstraintErr); ok {
			return ClientError{fmt.Errorf("api_id %d does not exist", constraintErr.value)}
		} else {
			return ServerError{err}
		}
	}

	return nil
}
