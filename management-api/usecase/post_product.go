package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"log"
)

func PostProduct(ctx context.Context, req *model.PostProductReq) error {
	// TODO: swaggerのパースとdynamoDBへの登録や、base_pathの取得
	dbParam := req.DBParam("/todo")
	log.Println(dbParam)
	if err := db.postProduct(ctx, &dbParam); err != nil {
		log.Printf("db insert api product error: %v", err)
		return ServerError{err}
	}
	return nil
}
