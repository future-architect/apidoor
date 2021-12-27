package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/future-architect/apidoor/gateway/model"
	"github.com/guregu/dynamo"
	"log"
	"os"
)

type APIForwarding struct {
	APIKey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}

func (af APIForwarding) Field() model.Field {
	return model.Field{
		Template: model.NewURITemplate(af.Path),
		Path:     model.NewURITemplate(af.ForwardURL),
		Num:      5,
		Max:      10,
	}
}

type DataSource struct {
	client             *dynamo.DB
	apiForwardingTable string
}

func New() *DataSource {
	apiForwardingTable := os.Getenv("DYNAMO_TABLE_API_FORWARDING")
	if apiForwardingTable == "" {
		log.Fatal("missing DYNAMO_TABLE_API_FORWARDING env")
	}

	dbEndpoint := os.Getenv("DYNAMO_ENDPOINT")
	if dbEndpoint != "" {
		return &DataSource{
			client: dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
			}))),
			apiForwardingTable: apiForwardingTable,
		}
	}

	return &DataSource{
		client:             dynamo.New(session.Must(session.NewSession())),
		apiForwardingTable: apiForwardingTable,
	}

}

func (dd DataSource) GetFields(ctx context.Context, key string) (model.Fields, error) {
	var resp []*APIForwarding
	err := dd.client.Table(dd.apiForwardingTable).
		Get("api_key", key).
		AllWithContext(ctx, &resp)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, model.ErrUnauthorizedRequest
		}
		return nil, &model.MyError{Message: fmt.Sprintf("internal server error: %v", err)}
	}

	fields := make([]model.Field, 0, len(resp))
	for _, forwarding := range resp {
		fields = append(fields, forwarding.Field())
	}

	return fields, nil
}
