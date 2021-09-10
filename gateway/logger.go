package gateway

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"time"
)

func UpdateLog(key, path string, r *http.Request) {
	// open log file
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0200)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	// check log pattern
	pattern := os.Getenv("LOG_PATTERN")
	schema, err := LogPatternParser(pattern)
	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: allow user to change order
	// make record
	record := []string{time.Now().Format(time.RFC3339), key, path}
	for _, value := range schema {
		record = append(record, r.Header.Get(value))
	}

	// write to log file
	writer := csv.NewWriter(file)
	writer.Write(record)
	writer.Flush()
}
