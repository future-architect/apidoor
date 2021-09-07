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
	//region := os.Getenv("AWS_REGION")

	if dbEndpoint == "" {
		// TODO: 本番環境向けの設定
		log.Fatal("missing DYNAMO_TABLE_API_FORWARDING env")
		return nil
	} else {
		return &DynamoDB{
			client: dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
			}))),
			apiForwardingTable: apiForwardingTable,
		}
	}
}

func (dd DynamoDB) GetFields(ctx context.Context, key string) (Fields, error) {
	var fields []Field

	var forwardings []*APIForwarding
	err := dd.client.Table(dd.apiForwardingTable).
		Get("api_key", key).
		AllWithContext(ctx, &forwardings)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, ErrUnauthorizedRequest
		} else {
			return nil, &MyError{Message: fmt.Sprintf("internal server error: %v", err)}
		}
	}

	for _, forwarding := range forwardings {
		u := NewURITemplate(forwarding.Path)
		v := NewURITemplate(forwarding.ForwardURL)
		fields = append(fields, Field{
			Template: *u,
			Path:     *v,
			Num:      5,
			Max:      10,
		})
	}

	return fields, nil
}
