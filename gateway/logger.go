package gateway

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
	"time"
)

type Log struct {
	sync.Mutex
	Data map[string]map[string]int
}

var TmpLog = Log{
	Data: make(map[string]map[string]int),
}

func UpdateLog(key, path string) {
	file, err := os.OpenFile("./log/log.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{time.Now().String(), key, path})
	writer.Flush()
}
