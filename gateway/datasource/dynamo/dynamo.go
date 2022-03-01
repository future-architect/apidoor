package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/future-architect/apidoor/gateway/datasource"
	"github.com/future-architect/apidoor/gateway/model"
	"github.com/guregu/dynamo"
	"log"
	"os"
)

type APIRouting struct {
	APIKey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}

type DataSource struct {
	client          *dynamo.DB
	apiRoutingTable string
	accessKeyTable  string
}

func New() *DataSource {
	apiRoutingTable := os.Getenv("DYNAMO_TABLE_API_ROUTING")
	if apiRoutingTable == "" {
		log.Fatal("missing DYNAMO_TABLE_API_ROUTING env")
	}
	//TODO: env

	dbEndpoint := os.Getenv("DYNAMO_DATA_SOURCE_ENDPOINT")
	if dbEndpoint != "" {
		return &DataSource{
			client: dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
			}))),
			apiRoutingTable: apiRoutingTable,
		}
	}

	return &DataSource{
		client:          dynamo.New(session.Must(session.NewSession())),
		apiRoutingTable: apiRoutingTable,
	}
}

func (dd DataSource) GetFields(ctx context.Context, key string) (model.Fields, error) {
	var routingList []*APIRouting
	err := dd.client.Table(dd.apiRoutingTable).
		Get("api_key", key).
		AllWithContext(ctx, &routingList)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, model.ErrUnauthorizedRequest
		}
		return nil, &model.MyError{Message: fmt.Sprintf("internal server error: %v", err)}
	}

	fields := make([]model.Field, 0, len(routingList))
	for _, routing := range routingList {
		field, err := datasource.CreateField(ctx, routing.APIKey, routing.Path, routing.ForwardURL)
		if err != nil {
			return nil, fmt.Errorf("fetch field, key = %v, hk = %v, forwardURL = %v, error: %w",
				routing.APIKey, routing.Path, routing.ForwardURL, err)
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func (dd DataSource) GetAccessTokens(ctx context.Context, apikey, templatePath string) (*model.AccessTokens, error) {
	key := fmt.Sprintf("%s#%s", apikey, templatePath)
	var tokens model.AccessTokens
	err := dd.client.Table(dd.accessKeyTable).
		Get("key", key).
		OneWithContext(ctx, &tokens)
	if err != nil && err != dynamo.ErrNotFound {
		return nil, &model.MyError{Message: fmt.Sprintf("get access tokens db error: %v", err)}
	}
	return &tokens, nil
}
