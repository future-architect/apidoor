package dblogger_test

import (
	"database/sql"
	"dblogger"
	"encoding/csv"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

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
		"2021-06-30 17:54:07.2278208 +0900 JST m=+0.032315701",
		"key",
		"path",
	})
	writer.Write([]string{
		"2021-06-30 17:54:07.2290175 +0900 JST m=+0.033512401",
		"key",
		"path",
	})
	writer.Flush()
	dblogger.PushLog()

	var n int
	if err := db.QueryRow("SELECT num from apilog WHERE apikey='key' AND apipath='path'").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("unexpected count %d, expected 2", n)
	}

	writer.Write([]string{
		"2021-07-01 10:28:32.9216589 +0900 JST m=+0.099747001",
		"key",
		"path",
	})
	writer.Flush()
	dblogger.PushLog()

	if err := db.QueryRow("SELECT num from apilog WHERE apikey='key' AND apipath='path'").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("unexpected count %d, expected 3", n)
	}

	if _, err := db.Exec("DELETE FROM apilog WHERE apikey='key' AND apipath='path'"); err != nil {
		t.Fatal(err)
	}
}
