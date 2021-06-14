package gateway

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Log map[string]map[string]int

var TmpLog = make(Log)

func UpdateLog(key, path string) {
	if _, ok := TmpLog[key][path]; !ok {
		TmpLog[key] = make(map[string]int)
		TmpLog[key][path] = 1
	} else {
		TmpLog[key][path]++
	}
}

func PushLog() {
	db, err := sql.Open(os.Getenv("DATABASE_DRIVER"),
		"host="+os.Getenv("DATABASE_HOST")+" "+
			"port="+os.Getenv("DATABASE_PORT")+" "+
			"user="+os.Getenv("DATABASE_USER")+" "+
			"password="+os.Getenv("DATABASE_PASSWORD")+" "+
			"dbname="+os.Getenv("DATABASE_NAME")+" "+
			"sslmode="+os.Getenv("DATABASE_SSLMODE"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for keyk, keyv := range TmpLog {
		for fieldk, fieldv := range keyv {
			tmp := struct {
				apikey string
				apilog string
				num    int
			}{
				apikey: "",
				apilog: "",
				num:    0,
			}

			switch err := db.QueryRow("SELECT * FROM apilog WHERE apikey=$1 AND apipath=$2", keyk, fieldk).Scan(&tmp.apikey, &tmp.apilog, &tmp.num); err {
			case sql.ErrNoRows:
				if _, err := db.Exec("INSERT INTO apilog(apikey, apipath, num) VALUES($1, $2, 1)", keyk, fieldk); err != nil {
					log.Fatal(err)
				}
			case nil:
				if _, err := db.Exec("UPDATE apilog SET num=num+$1 WHERE apikey=$2 AND apipath=$3", fieldv, keyk, fieldk); err != nil {
					log.Fatal(err)
				}
			default:
				log.Fatal(err)
			}
		}
	}

	TmpLog = make(Log)
}
