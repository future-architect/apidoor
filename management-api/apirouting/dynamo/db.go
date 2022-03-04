package dynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/guregu/dynamo"
	"log"
	"os"
)

type APIRouting struct {
	client           *dynamo.DB
	apiRoutingTable  string
	accessTokenTable string
}

func New() *APIRouting {
	apiRoutingTable := os.Getenv("DYNAMO_TABLE_API_ROUTING")
	if apiRoutingTable == "" {
		log.Fatal("missing DYNAMO_TABLE_API_ROUTING env")
	}
	accessTokenTable := os.Getenv("DYNAMO_TABLE_ACCESS_TOKEN")
	if apiRoutingTable == "" {
		log.Fatal("missing DYNAMO_TABLE_API_TOKEN env")
	}

	var client *dynamo.DB

	dbEndpoint := os.Getenv("DYNAMO_ENDPOINT")
	if dbEndpoint != "" {
		client = dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           "local",
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
		})))
	} else {
		client = dynamo.New(session.Must(session.NewSession()))
	}

	return &APIRouting{
		client:           client,
		apiRoutingTable:  apiRoutingTable,
		accessTokenTable: accessTokenTable,
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

func (ar APIRouting) CountRouting(ctx context.Context, apikey, path string) (int64, error) {
	return ar.client.Table(ar.apiRoutingTable).
		Get("api_key", apikey).
		Range("path", dynamo.Equal, path).
		CountWithContext(ctx)
}

func (ar APIRouting) PostAPIToken(ctx context.Context, req model.PostAPITokenReq) error {
	accessTokens := newAccessToken(req)
	return ar.client.Table(ar.accessTokenTable).
		Put(accessTokens).RunWithContext(ctx)
}

func (ar APIRouting) DeleteAPIToken(ctx context.Context, req model.DeleteAPITokenReq) error {
	key := fmt.Sprintf("%s#%s", req.APIKey, req.Path)
	return ar.client.Table(ar.accessTokenTable).
		Delete("key", key).
		RunWithContext(ctx)
}

type routing struct {
	Apikey     string `dynamo:"api_key"`
	Path       string `dynamo:"path"`
	ForwardURL string `dynamo:"forward_url"`
}

type accessTokens struct {
	Key          string              `dynamo:"key"` // <api_key>#<path>
	AccessTokens []model.AccessToken `dynamo:"tokens"`
}

func newAccessToken(req model.PostAPITokenReq) accessTokens {
	key := fmt.Sprintf("%s#%s", req.APIKey, req.Path)
	return accessTokens{
		Key:          key,
		AccessTokens: req.AccessTokens,
	}
}
