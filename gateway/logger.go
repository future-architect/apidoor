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

	// make record
	record := []string{}
	for _, value := range schema {
		switch value {
		case "time":
			record = append(record, time.Now().Format(time.RFC3339))
		case "key":
			record = append(record, key)
		case "path":
			record = append(record, path)
		default:
			record = append(record, r.Header.Get(value))
		}
	}

	// write to log file
	writer := csv.NewWriter(file)
	writer.Write(record)
	writer.Flush()
}
