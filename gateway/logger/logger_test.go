package logger_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/future-architect/apidoor/gateway"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/guregu/dynamo"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

var (
	db             *dynamo.DB
	accessLogTable string
)

func init() {
	accessLogTable = os.Getenv("DYNAMO_TABLE_ACCESS_LOG")
	if accessLogTable == "" {
		log.Fatal("missing DYNAMO_TABLE_ACCESS_LOG env")
	}

	dbEndpoint := os.Getenv("DYNAMO_ACCESS_LOG_ENDPOINT")
	db = dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "local",
		Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
	})))

}

func calcBillingStatus(resp *http.Response) logger.BillingStatus {
	code := resp.StatusCode
	msd := code / 100
	if msd == 4 || msd == 5 {
		return logger.NotBilling
	}
	return logger.Billing
}

func TestUpdateLog(t *testing.T) {

	restore := flextime.Set(time.Date(2021, time.December, 27, 17, 1, 41, 0, time.UTC))
	defer restore()

	tests := []struct {
		name           string
		contractID     int
		key            string
		path           string
		header         map[string]string
		responseStatus int
		logPattern     string
		wantLog        string
	}{
		{
			name:           "billing status is billing",
			contractID:     0,
			key:            "key",
			path:           "path",
			responseStatus: http.StatusOK,
			logPattern:     "time,key,path,response_status,billing_status",
			wantLog:        "2021-12-27T17:01:41Z,key,path,200,billing\n",
		},
		{
			name:           "billing status is not billing",
			contractID:     0,
			key:            "key",
			path:           "path",
			responseStatus: http.StatusInternalServerError,
			logPattern:     "time,key,path,response_status,billing_status",
			wantLog:        "2021-12-27T17:01:41Z,key,path,500,not billing\n",
		},
		{
			name:       "write header values",
			contractID: 0,
			key:        "key",
			path:       "path",
			header: map[string]string{
				"TEST1": "header1",
				"TEST3": "header3",
			},
			responseStatus: http.StatusOK,
			logPattern:     "time,key,path,TEST1,TEST2,TEST3",
			wantLog:        "2021-12-27T17:01:41Z,key,path,header1,,header3\n",
		},
	}

	defer func() {
		os.Setenv("LOG_PATTERN", "")
		logger.LogOptionPattern = logger.DefaultLogPattern()
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv("LOG_PATTERN", tt.logPattern); err != nil {
				t.Errorf("set ptn env failed: %v", err)
				return
			}
			logger.LogOptionPattern = logger.DefaultLogPattern()

			buffer := &bytes.Buffer{}
			appender := logger.DefaultAppender{
				Writer: buffer,
			}

			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			for k, v := range tt.header {
				r.Header.Set(k, v)
			}
			resp := http.Response{
				StatusCode: tt.responseStatus,
			}

			appender.Do(tt.contractID, tt.key, tt.path, r, &resp, calcBillingStatus)

			want := strings.ReplaceAll(tt.wantLog, "\r\n", "\n")
			//want := tt.wantLog
			if out := buffer.String(); out != want {
				t.Errorf("result:\n%s\nwant:\n%s", out, want)
			}
		})
	}

}

func TestUpdateDBRoutine(t *testing.T) {
	gateway.Setup(t,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../../dynamo_table/access_log_table.json`,
	)
	t.Cleanup(func() {
		gateway.Teardown(t,
			`aws dynamodb --profile local --endpoint-url http://localhost:4566 delete-table --table access_log`,
		)
	})

	appender := logger.NewCSVAppender(csv.NewWriter(io.Discard))
	inputAccesses := []struct {
		contractID int
		key        string
		path       string
		r          *http.Request
		resp       http.Response
	}{
		{
			contractID: 0,
			key:        "api_key1",
			path:       "example.com/api1",
			r:          httptest.NewRequest(http.MethodGet, "localhost:3000/products", nil),
			resp:       http.Response{StatusCode: http.StatusOK},
		},
		{
			contractID: 1,
			key:        "api_key2",
			path:       "example.com/api2",
			r:          httptest.NewRequest(http.MethodGet, "localhost:3000/products", nil),
			resp:       http.Response{StatusCode: http.StatusInternalServerError},
		},
	}

	//start simulation
	routineKill := make(chan bool)
	routineFinish := make(chan bool)
	ctx := context.Background()

	go logger.UpdateDBRoutine(ctx, &appender, 5*time.Second, routineKill, routineFinish)
	defer logger.CleanupUpdateDBTask(routineKill, routineFinish)

	if err := appender.Do(inputAccesses[0].contractID, inputAccesses[0].key, inputAccesses[0].path, inputAccesses[0].r,
		&inputAccesses[0].resp, calcBillingStatus); err != nil {
		t.Errorf("append access log %+v failed: %v", inputAccesses[0], err)
	}

	testAccessLogDBResult(t, "the process has not put items", []logger.LogItem{})

	// updating the db occurs during the sleep
	time.Sleep(7 * time.Second)

	testAccessLogDBResult(t, "the process has put an item", []logger.LogItem{
		{
			ContractID:    inputAccesses[0].contractID,
			Key:           inputAccesses[0].key,
			Path:          inputAccesses[0].path,
			StatusCode:    http.StatusOK,
			BillingStatus: logger.Billing,
		},
	})

	if err := appender.Do(inputAccesses[1].contractID, inputAccesses[1].key, inputAccesses[1].path, inputAccesses[1].r,
		&inputAccesses[1].resp, calcBillingStatus); err != nil {
		t.Errorf("append access log %+v failed: %v", inputAccesses[1], err)
	}

	testAccessLogDBResult(t, "the process has not put the second item",
		[]logger.LogItem{
			{
				ContractID:    inputAccesses[0].contractID,
				Key:           inputAccesses[0].key,
				Path:          inputAccesses[0].path,
				StatusCode:    http.StatusOK,
				BillingStatus: logger.Billing,
			},
		})

	// updating the db occurs during the sleep
	time.Sleep(5 * time.Second)

	testAccessLogDBResult(t, "the process has put the second item and duplicate putting an item has not occurred",
		[]logger.LogItem{
			{
				ContractID:    inputAccesses[1].contractID,
				Key:           inputAccesses[1].key,
				Path:          inputAccesses[1].path,
				StatusCode:    http.StatusInternalServerError,
				BillingStatus: logger.NotBilling,
			},
			{
				ContractID:    inputAccesses[0].contractID,
				Key:           inputAccesses[0].key,
				Path:          inputAccesses[0].path,
				StatusCode:    http.StatusOK,
				BillingStatus: logger.Billing,
			},
		})
}

func testAccessLogDBResult(t *testing.T, name string, wantResult []logger.LogItem) {
	result := scanAccessLog(t)

	if len(result) < 1 {
		if len(wantResult) >= 1 {
			t.Errorf("check %v: scan access log result is empty, but want: %+v", name, wantResult)
		}
		return
	}
	if diff := cmp.Diff(result, wantResult, cmpopts.IgnoreFields(logger.LogItem{}, "TimeStamp")); diff != "" {
		t.Errorf("check %v: access log scan result differs:\n%v", name, diff)
	}
}

func scanAccessLog(t *testing.T) []logger.LogItem {
	result := make([]logger.LogItem, 0)
	if err := db.Table(accessLogTable).Scan().All(&result); err != nil {
		t.Errorf("scan access log failed: %v", err)
	}
	return result
}
