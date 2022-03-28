package usecase

import (
	"context"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/apirouting"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/swagger-parser"
	"log"
)

var Parser = swaggerparser.NewParser(swaggerparser.NewDefaultFetcher())

func PostProduct(ctx context.Context, req *model.PostProductReq) (*model.Product, error) {
	swaggerInfo, err := Parser.Parse(ctx, req.SwaggerURL)
	if err != nil {
		parseErr, _ := err.(swaggerparser.Error)
		log.Printf("failed to fetch and parse swagger definition file: %v", err)
		switch parseErr.ErrorType {
		case swaggerparser.FetchServerError, swaggerparser.FileParseError, swaggerparser.FormatError:
			return nil, ClientError{fmt.Errorf("failed to fetch and parse swagger definition file: %w", err)}
		default:
			return nil, ServerError{err}
		}
	}

	dbParam := req.DBParam(swaggerInfo.PathBase)

	product, err := db.postProduct(ctx, &dbParam)
	if err != nil {
		log.Printf("db insert api product error: %v", err)
		return nil, ServerError{err}
	}

	if err = apirouting.ApiDBDriver.PostSwagger(ctx, product.ID, swaggerInfo); err != nil {
		log.Printf("db insert swagger error: %v", err)
		log.Printf("delete product, id = %d", product.ID)

		if err = db.deleteProduct(ctx, product.ID); err != nil {
			log.Printf("db delete product error: %v", err)
			return nil, ServerError{err}
		}
	}

	return product, nil
}
