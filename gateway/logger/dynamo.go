package logger

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"log"
	"os"
	"time"
)

type accessLogDB struct {
	client         *dynamo.DB
	accessLogTable string
}

var (
	db accessLogDB
)

func init() {
	accessLogTable := os.Getenv("DYNAMO_TABLE_ACCESS_LOG")
	if accessLogTable == "" {
		log.Fatal("missing DYNAMO_TABLE_ACCESS_LOG env")
	}

	dbEndpoint := os.Getenv("DYNAMO_ACCESS_LOG_ENDPOINT")
	if dbEndpoint != "" {
		db = accessLogDB{
			client: dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
				Profile:           "local",
				SharedConfigState: session.SharedConfigEnable,
				Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
			}))),
			accessLogTable: accessLogTable,
		}
	} else {
		db = accessLogDB{
			client:         dynamo.New(session.Must(session.NewSession())),
			accessLogTable: accessLogTable,
		}
	}
}

func (ad accessLogDB) postAccessLogDB(ctx context.Context, item LogItem) error {
	return ad.client.Table(ad.accessLogTable).
		Put(item).RunWithContext(ctx)
}

func (ad accessLogDB) countBillingAccessLogDB(ctx context.Context, contractID int, path string, startAt time.Time) (int64, error) {
	return ad.client.Table(ad.accessLogTable).
		Get("contract_id", contractID).
		Range("timestamp", dynamo.GreaterOrEqual, startAt).
		Filter("'path' = ?", path).
		Filter("'billing_status' = ?", Billing).
		CountWithContext(ctx)
}
