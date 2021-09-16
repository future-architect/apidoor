package dblogger_test

import (
	"database/sql"
	"dblogger"
	"encoding/csv"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

type testLog = struct {
	date   string
	key    string
	path   string
	custom string
}

var testLogData = testLog{
	date:   "2021-07-01T14:01:46+09:00",
	key:    "key",
	path:   "path",
	custom: `{"key": "value"}`,
}

func TestPushLog(t *testing.T) {
	// set up database
	db, err := sql.Open(os.Getenv("DATABASE_DRIVER"),
		"host="+os.Getenv("DATABASE_HOST")+" "+
			"port="+os.Getenv("DATABASE_PORT")+" "+
			"user="+os.Getenv("DATABASE_USER")+" "+
			"password="+os.Getenv("DATABASE_PASSWORD")+" "+
			"dbname="+os.Getenv("DATABASE_NAME")+" "+
			"sslmode="+os.Getenv("DATABASE_SSLMODE"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// clean up database for test
	if _, err := db.Exec("DELETE FROM log_list WHERE api_key='key' AND api_path='path'"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM log_list WHERE api_key='key' AND api_path='path'")
	})

	// open log file
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}

	// execute PushLog()
	writer := csv.NewWriter(file)
	writer.Write([]string{
		testLogData.date,
		testLogData.key,
		testLogData.path,
		testLogData.custom,
	})
	writer.Flush()
	dblogger.PushLog()

	// check if log is written to database correctly
	row := testLog{}
	if err := db.QueryRow("SELECT * FROM log_list WHERE api_key='key'").Scan(&row.date, &row.key, &row.path, &row.custom); err != nil {
		t.Fatal(err)
	}

	t1, err := time.Parse(time.RFC3339, testLogData.date)
	if err != nil {
		t.Fatal(err)
	}
	t2, err := time.Parse(time.RFC3339, row.date)
	if err != nil {
		t.Fatal(err)
	}

	if !t1.Equal(t2) {
		t.Fatalf("unexpected date %s, expected %s", t1.String(), t2.String())
	} else if row.key != testLogData.key {
		t.Fatalf("unexpected key %s, expected %s", row.key, testLogData.key)
	} else if row.path != testLogData.path {
		t.Fatalf("unexpected path %s, expected %s", row.path, testLogData.path)
	} else if row.custom != testLogData.custom {
		t.Fatalf("unexpected custom data %s, expected %s", row.custom, testLogData.custom)
	}
}
