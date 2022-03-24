package usecase

import (
	"context"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/swagger-parser"
	"log"
)

func PostProduct(ctx context.Context, req *model.PostProductReq) error {
	// TODO: swaggerのパースとdynamoDBへの登録や、base_pathの取得
	dbParam := req.DBParam("/todo")
	_ = swagger_parser.NewParser(swagger_parser.NewDefaultFetcher())
	if err := db.postProduct(ctx, &dbParam); err != nil {
		log.Printf("db insert api product error: %v", err)
		return ServerError{err}
	}

	return nil
}
