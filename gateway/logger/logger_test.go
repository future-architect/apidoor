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

func TestUpdateLog(t *testing.T) {

	restore := flextime.Set(time.Date(2021, time.December, 27, 17, 1, 41, 0, time.UTC))
	defer restore()

	r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Set("TEST1", "header1")
	r.Header.Set("TEST2", "header2")
	logger.LogOptionPattern = []logger.LogOption{
		logger.WithTime(),
		logger.WithKey(),
		logger.WithPath(),
		logger.HeaderElement("TEST1"),
		logger.HeaderElement("TEST2"),
	}

	buffer := &bytes.Buffer{}
	appender := logger.DefaultAppender{
		Writer: buffer,
	}
	for i := 0; i < 2; i++ {
		appender.Do("key", "path", r)
	}

	want := `2021-12-27T17:01:41Z,key,path,header1,header2
2021-12-27T17:01:41Z,key,path,header1,header2
`
	want = strings.ReplaceAll(want, "\r\n", "\n")

	if out := buffer.String(); out != want {
		t.Fatalf("result:\n%s\nwant:\n%s", out, want)
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
		key  string
		path string
		r    *http.Request
	}{
		{
			key:  "api_key1",
			path: "example.com/api1",
			r:    httptest.NewRequest(http.MethodGet, "localhost:3000/products", nil),
		},
		{
			key:  "api_key2",
			path: "example.com/api2",
			r:    httptest.NewRequest(http.MethodGet, "localhost:3000/products", nil),
		},
	}

	//start simulation
	routineKill := make(chan bool)
	routineFinish := make(chan bool)
	ctx := context.Background()

	go logger.UpdateDBRoutine(ctx, &appender, 5*time.Second, routineKill, routineFinish)
	defer logger.CleanupUpdateDBTask(routineKill, routineFinish)

	if err := appender.Do(inputAccesses[0].key, inputAccesses[0].path, inputAccesses[0].r); err != nil {
		t.Errorf("append access log %+v failed: %v", inputAccesses[0], err)
	}

	testAccessLogDBResult(t, "the process has not put items", []logger.LogItem{})

	// updating the db occurs during the sleep
	time.Sleep(7 * time.Second)

	testAccessLogDBResult(t, "the process has put an item", []logger.LogItem{
		{
			Key:  inputAccesses[0].key,
			Path: inputAccesses[0].path,
		},
	})

	if err := appender.Do(inputAccesses[1].key, inputAccesses[1].path, inputAccesses[1].r); err != nil {
		t.Errorf("append access log %+v failed: %v", inputAccesses[1], err)
	}

	testAccessLogDBResult(t, "the process has not put the second item",
		[]logger.LogItem{
			{
				Key:  inputAccesses[0].key,
				Path: inputAccesses[0].path,
			},
		})

	// updating the db occurs during the sleep
	time.Sleep(5 * time.Second)

	testAccessLogDBResult(t, "the process has put the second item and duplicate putting an item has not occurred",
		[]logger.LogItem{
			{
				Key:  inputAccesses[1].key,
				Path: inputAccesses[1].path,
			},
			{
				Key:  inputAccesses[0].key,
				Path: inputAccesses[0].path,
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
