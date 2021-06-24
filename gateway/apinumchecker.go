package gateway

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

func ApiNumChecker(key, path string) error {
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

	var max int
	switch err := db.QueryRow("SELECT apimaxnum FROM apilimit WHERE apikey=$1 AND apipath=$2", key, path).Scan(&max); err {
	case sql.ErrNoRows:
		return nil
	case nil:
		if k, ok := TmpLog.Data[key][path]; ok {
			n += k
		}
		if n >= max {
			return errors.New("limit exceeded")
		}
		return nil
	default:
		return errors.New("error occurs in database")
	}
}
