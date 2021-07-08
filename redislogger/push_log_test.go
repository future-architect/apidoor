package redislogger_test

import (
	"context"
	"encoding/csv"
	"os"
	"redislogger"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var now = time.Now()

var testdata = []struct {
	date string
	key  string
	path string
}{
	// data not counted by logger
	{
		date: now.Add(-2 * time.Minute).Format(time.RFC3339),
		key:  "key",
		path: "path",
	},
	// data counted by logger
	{
		date: now.Add(-2 * time.Second).Format(time.RFC3339),
		key:  "key",
		path: "path",
	},
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

func TestPushLog(t *testing.T) {
	file, err := os.OpenFile("./log/log.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}

	writer := csv.NewWriter(file)
	for _, tt := range testdata {
		writer.Write([]string{
			tt.date,
			tt.key,
			tt.path,
		})
	}
	writer.Flush()

	ctx := context.Background()
	rdb.HDel(ctx, "key", "path")

	redislogger.PushLog()

	n, err := strconv.Atoi(rdb.HGet(ctx, "key", "path").Val())
	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Fatalf("unexpected count %d, expected %d", n, 1)
	}

	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}
	rdb.HDel(ctx, "key", "path")
}
