package gateway

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

func ApiLimitChecker(key, path string) error {
	TmpLog.Lock()
	defer TmpLog.Unlock()

	db, err := sql.Open(os.Getenv("DATABASE_DRIVER"),
		"host="+os.Getenv("DATABASE_HOST")+" "+
			"port="+os.Getenv("DATABASE_PORT")+" "+
			"user="+os.Getenv("DATABASE_USER")+" "+
			"password="+os.Getenv("DATABASE_PASSWORD")+" "+
			"dbname="+os.Getenv("DATABASE_NAME")+" "+
			"sslmode="+os.Getenv("DATABASE_SSLMODE"))
	if err != nil {
		return errors.New("error occurs in database")
	}
	defer db.Close()

	var n int
	if err := db.QueryRow("SELECT num FROM apilog WHERE apikey=$1 AND apipath=$2", key, path).Scan(&n); err != nil && err != sql.ErrNoRows {
		return errors.New("error occurs in database")
	}

	k := TmpLog.Data[key][path]
	for _, field := range APIData[key] {
		if field.Template.JoinPath() == path {
			if n+k >= field.Max {
				return errors.New("limit exceeded")
			}
		}
	}

	return nil
}
