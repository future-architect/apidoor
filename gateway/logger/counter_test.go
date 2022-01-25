package logger_test

import (
	"context"
	"github.com/Songmu/flextime"
	"github.com/future-architect/apidoor/gateway"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/future-architect/apidoor/gateway/model"
	"testing"
	"time"
)

var (
	testCounter = logger.APICallCounter{}
)

func TestAPICallCounter_GetCount(t *testing.T) {
	gateway.Setup(t,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../db_table/access_log_table.json`,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 batch-write-item --request-items file://./testdata/get_counter_items.json`,
	)
	t.Cleanup(func() {
		gateway.Teardown(t,
			`aws dynamodb --profile local --endpoint-url http://localhost:4566 delete-table --table access_log`,
		)
	})
	restore := flextime.Fix(time.Date(2020, 12, 15, 13, 0, 0, 0, time.UTC))
	defer restore()

	logger.DefaultCountSpanDays = 30

	tests := []struct {
		name      string
		apikey    string
		path      model.URITemplate
		wantCount int
	}{
		{
			name:      "count api calls correctly",
			apikey:    "key1",
			path:      model.NewURITemplate("api/test"),
			wantCount: 2,
		},
		{
			name:      "count result is zero",
			apikey:    "key9",
			path:      model.NewURITemplate("api/test"),
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := testCounter.GetCount(context.Background(), tt.apikey, tt.path)
			if err != nil {
				t.Errorf("get count error: %v", err)
			}

			if count != tt.wantCount {
				t.Errorf("count result is wrong, want %d, got %d", tt.wantCount, count)
			}
		})
	}

}

func TestAPICallCounter_GetCountWithCache(t *testing.T) {
	gateway.Setup(t,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../db_table/access_log_table.json`,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 batch-write-item --request-items file://./testdata/get_counter_items.json`,
	)
	t.Cleanup(func() {
		gateway.Teardown(t,
			`aws dynamodb --profile local --endpoint-url http://localhost:4566 delete-table --table access_log`,
		)
	})

	startTime := time.Date(2020, 12, 15, 13, 0, 0, 0, time.UTC)

	restore := flextime.Fix(startTime)
	defer restore()

	logger.DefaultCountSpanDays = 30
	logger.DefaultCountValidSpan = 30 * time.Second

	apikey := "key1"
	pathStr := "api/test"
	path := model.NewURITemplate(pathStr)

	// get count correctly
	testGetCount(t, apikey, path, 2)

	flextime.Fix(startTime.Add(5 * time.Second))
	putAccessLog(t, logger.LogItem{
		Key:       apikey,
		TimeStamp: startTime.Add(5 * time.Second).Format(time.RFC3339),
		Path:      pathStr,
	})

	// the result is not changed, because the cache is valid
	testGetCount(t, apikey, path, 2)

	flextime.Fix(startTime.Add(35 * time.Second))
	// the result is updated, because the cache is invalid
	testGetCount(t, apikey, path, 3)

}

func testGetCount(t *testing.T, key string, path model.URITemplate, wantCount int) {
	count, err := testCounter.GetCount(context.Background(), key, path)
	if err != nil {
		t.Errorf("get count error: %v", err)
	}
	if count != wantCount {
		t.Errorf("count result is wrong, want %d, got %d", wantCount, count)
	}
}

func putAccessLog(t *testing.T, item logger.LogItem) {
	if err := db.Table(accessLogTable).Put(item).Run(); err != nil {
		t.Errorf("put access log failed: %v", err)
	}
}
