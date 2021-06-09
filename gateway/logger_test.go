package gateway_test

import (
	"database/sql"
	"gateway"
	"testing"

	_ "github.com/lib/pq"
)

func TestUpdateLog(t *testing.T) {
	for i := 1; i <= 2; i++ {
		gateway.UpdateLog("key", "path")
		if gateway.TmpLog["key"]["path"] != i {
			t.Fatalf("unexpected TmpLog[key][path]: %d, expected %d", gateway.TmpLog["key"]["path"], i)
		}
	}

	if _, ok := gateway.TmpLog["key"]["unusedpath"]; ok {
		t.Fatal("unexpected field in data")
	}
}

func TestPushLog(t *testing.T) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5555 user=root password=password dbname=root sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tmp := struct {
		apikey string
		apilog string
		num    int
	}{
		apikey: "",
		apilog: "",
		num:    0,
	}

	for i := 1; i <= 2; i++ {
		gateway.UpdateLog("key", "path")
		gateway.PushLog()
		switch err := db.QueryRow("SELECT * FROM apilog WHERE apikey=$1 AND apipath=$2", "key", "path").Scan(&tmp.apikey, &tmp.apilog, &tmp.num); err {
		case sql.ErrNoRows:
			t.Fatal("there is no expected data in database")
		case nil:
			if tmp.num != i {
				t.Fatalf("unexpected API count %d, expected %d", tmp.num, i)
			}
		default:
			t.Fatal(err)
		}
	}
}
