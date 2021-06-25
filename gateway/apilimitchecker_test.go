package gateway_test

import (
	"database/sql"
	"errors"
	"gateway"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

type apinumtest struct {
	data  int
	local int
	max   int
	err   error
}

var apinumdata = []apinumtest{
	// valid request
	{
		data:  1,
		local: 1,
		max:   4,
		err:   nil,
	},
	// valid request(boundary)
	{
		data:  1,
		local: 1,
		max:   3,
		err:   nil,
	},
	// valid request(permit unlimited call)
	{
		data:  1,
		local: 1,
		max:   -1,
		err:   nil,
	},
	// invalid request
	{
		data:  2,
		local: 3,
		max:   4,
		err:   errors.New("limit exceeded"),
	},
	// invalid request(boundary)
	{
		data:  2,
		local: 2,
		max:   4,
		err:   errors.New("limit exceeded"),
	},
}

func TestAPINumChecker(t *testing.T) {
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

	for i, tt := range apinumdata {
		if _, err := db.Exec("INSERT INTO apilog(apikey, apipath, num) VALUES('key', 'path', $1)", tt.data); err != nil {
			t.Fatalf("case %d: error occurs in database, %v", i, err)
		}

		gateway.TmpLog.Data = make(map[string]map[string]int)
		gateway.TmpLog.Data["key"] = make(map[string]int)
		gateway.TmpLog.Data["key"]["path"] = tt.local

		if tt.max >= 0 {
			gateway.APIData["key"] = []gateway.Field{
				{
					Template: *gateway.NewURITemplate("/path"),
					Path:     *gateway.NewURITemplate("/path"),
					Max:      tt.max,
				},
			}
		}

		switch err := gateway.ApiLimitChecker("key", "path"); err {
		case nil:
			if tt.err != nil {
				t.Fatalf("case %d: expected %v, get %v", i, tt.err, err)
			}
		default:
			if tt.err == nil {
				t.Fatalf("case %d: expected %v, get %v", i, tt.err, err)
			}
			if err.Error() != tt.err.Error() {
				t.Fatalf("case %d: expected %v, get %v", i, tt.err, err)
			}
		}

		gateway.TmpLog.Data = make(map[string]map[string]int)
		db.Exec("DELETE FROM apilog WHERE apikey='key'")
		db.Exec("DELETE FROM apilimit WHERE apikey='key'")
	}
}
