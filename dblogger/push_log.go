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
	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		date := line[0]
		key := line[1]
		path := line[2]
		if _, err := db.Exec("INSERT INTO apilog(rundate, apikey, apipath) VALUES($1, $2, $3)", date, key, path); err != nil {
			log.Fatal(err)
		}
	}

	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}
}
