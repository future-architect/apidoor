package gateway

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"log"
	"os"
)

type DynamoAPIForwarding struct {
	APIKey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}

func (af DynamoAPIForwarding) Field() Field {
	return Field{
		Template: *NewURITemplate(af.Path),
		Path:     *NewURITemplate(af.ForwardURL),
		Num:      5,
		Max:      10,
	}
}

type DynamoDB struct {
	client             *dynamo.DB
	apiForwardingTable string
}

func NewDynamoDB() *DynamoDB {
	apiForwardingTable := os.Getenv("DYNAMO_TABLE_API_FORWARDING")
	if apiForwardingTable == "" {
		log.Fatal("missing DYNAMO_TABLE_API_FORWARDING env")
	}

	dbEndpoint := os.Getenv("DYNAMO_ENDPOINT")
	if dbEndpoint != "" {
		return &DynamoDB{
			client: dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
			}))),
			apiForwardingTable: apiForwardingTable,
		}
	}

	return &DynamoDB{
		client:             dynamo.New(session.Must(session.NewSession())),
		apiForwardingTable: apiForwardingTable,
	}

}

func (dd DynamoDB) GetFields(ctx context.Context, key string) (Fields, error) {
	var resp []*DynamoAPIForwarding
	err := dd.client.Table(dd.apiForwardingTable).
		Get("api_key", key).
		AllWithContext(ctx, &resp)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, ErrUnauthorizedRequest
		} else {
			return nil, &MyError{Message: fmt.Sprintf("internal server error: %v", err)}
		}
	}

	fields := make([]Field, 0, len(resp))
	for _, forwarding := range resp {
		fields = append(fields, forwarding.Field())
	}

	return fields, nil
}
