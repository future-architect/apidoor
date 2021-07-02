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

var testdata = struct {
	date string
	key  string
	path string
}{
	date: "2021-07-01T14:01:46+09:00",
	key:  "key",
	path: "path",
}

func TestPushLog(t *testing.T) {
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
	if _, err := db.Exec("DELETE FROM apilog WHERE apikey='key' AND apipath='path'"); err != nil {
		t.Fatal(err)
	}

	file, err := os.OpenFile("./log/log.csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}

	writer := csv.NewWriter(file)
	writer.Write([]string{
		testdata.date,
		testdata.key,
		testdata.path,
	})
	writer.Flush()
	dblogger.PushLog()

	row := struct {
		date string
		key  string
		path string
	}{}
	if err := db.QueryRow("SELECT * FROM apilog WHERE apikey='key' AND apipath='path'").Scan(&row.date, &row.key, &row.path); err != nil {
		t.Fatal(err)
	}

	d, err := time.Parse(time.RFC3339, testdata.date)
	if err != nil {
		t.Fatal(err)
	}
	row.date = d.Format(time.RFC3339)

	if row.date != testdata.date {
		t.Fatalf("unexpected date %s, expected %s", row.date, testdata.date)
	} else if row.key != testdata.key {
		t.Fatalf("unexpected key %s, expected %s", row.key, testdata.key)
	} else if row.path != testdata.path {
		t.Fatalf("unexpected path %s, expected %s", row.path, testdata.path)
	}

	if _, err := db.Exec("DELETE FROM apilog WHERE apikey='key' AND apipath='path'"); err != nil {
		t.Fatal(err)
	}
}
