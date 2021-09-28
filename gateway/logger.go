package gateway

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
)

func UpdateLog(key, path string, r *http.Request) {
	// open log file
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0200)
	if err != nil {
		log.Fatal(err.Error())
	}

	// make record
	record := []string{}
	for _, logOption := range LogOptionPattern {
		logOption(&record, key, path, r)
	}

	// write to log file
	writer := csv.NewWriter(file)
	writer.Write(record)
	writer.Flush()
}
