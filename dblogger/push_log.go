package dblogger

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func PushLog() {
	// set up database
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

	// open log file
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// write log to database
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
		custom := line[3]
		if _, err := db.Exec("INSERT INTO log_list(run_date, api_key, api_path, custom_log) VALUES($1, $2, $3, $4)", date, key, path, custom); err != nil {
			log.Fatal(err)
		}
	}

	// initialize log file
	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}
}
