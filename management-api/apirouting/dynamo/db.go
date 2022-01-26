package dynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"log"
	"os"
)

type APIRouting struct {
	client          *dynamo.DB
	apiRoutingTable string
}

func New() *APIRouting {
	apiRoutingTable := os.Getenv("DYNAMO_TABLE_API_ROUTING")
	if apiRoutingTable == "" {
		log.Fatal("missing DYNAMO_TABLE_API_ROUTING env")
	}

	dbEndpoint := os.Getenv("DYNAMO_ENDPOINT")
	if dbEndpoint != "" {
		return &APIRouting{
			client: dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
			}))),
			apiRoutingTable: apiRoutingTable,
		}
	}

	return &APIRouting{
		client:          dynamo.New(session.Must(session.NewSession())),
		apiRoutingTable: apiRoutingTable,
	}
}

func (ar APIRouting) PostRouting(ctx context.Context, apikey, path, forwardURL string) error {
	routing := routing{
		Apikey:     apikey,
		Path:       path,
		ForwardURL: forwardURL,
	}
	return ar.client.Table(ar.apiRoutingTable).
		Put(routing).RunWithContext(ctx)
}

type routing struct {
	Apikey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}
