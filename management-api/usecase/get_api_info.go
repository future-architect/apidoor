package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func GetAPIInfo(ctx context.Context) ([]model.APIInfo, error) {
	list, err := db.getAPIInfo(ctx)
	if err != nil {
		log.Printf("execute get apiinfo from db error: %v", err)
		return nil, err
	}
	return list, nil
}
