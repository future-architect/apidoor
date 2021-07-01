package dblogger

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	_ "github.com/lib/pq"
)

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

	file, err := os.OpenFile("./log/log.csv", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	count := make(map[string]map[string]int)
	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		key := line[1]
		path := line[2]
		if _, ok := count[key][path]; ok {
			count[key][path]++
		} else {
			if _, ok := count[key]; !ok {
				count[key] = make(map[string]int)
			}
			count[key][path] = 1
		}
	}

	for k, v := range count {
		for p, n := range v {
			var isExist bool
			if err := db.QueryRow("SELECT EXISTS(SELECT * FROM apilog WHERE apikey=$1 AND apipath=$2)", k, p).Scan(&isExist); err != nil {
				log.Fatal(err)
			}

			if isExist {
				if _, err := db.Exec("UPDATE apilog SET num=num+$1 WHERE apikey=$2 AND apipath=$3", n, k, p); err != nil {
					log.Fatal(err)
				}
			} else {
				if _, err := db.Exec("INSERT INTO apilog(apikey, apipath, num) VALUES($1, $2, $3)", k, p, n); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}
}
