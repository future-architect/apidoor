package gateway

import (
	"net/http"
	"time"
)

type LogElement string

func (element LogElement) append(record []string, key, path string, r *http.Request) []string {
	switch element {
	case "time":
		record = append(record, time.Now().Format(time.RFC3339))
	case "key":
		record = append(record, key)
	case "path":
		record = append(record, path)
	default:
		record = append(record, r.Header.Get(string(element)))
	}

	return record
}
