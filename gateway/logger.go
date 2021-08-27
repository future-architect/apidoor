package gateway

import (
	"encoding/csv"
	"log"
	"os"
	"time"
)

func UpdateLog(key, path string) {
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0200)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{time.Now().Format(time.RFC3339), key, path})
	writer.Flush()
}
