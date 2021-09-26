package gateway

import (
	"fmt"
	"net/http"
)

func UpdateLog(key, path string, r *http.Request) {
	// make record
	record := []string{}
	for _, parsedElement := range LogPattern {
		record = parsedElement.append(record, key, path, r)
	}
	// output
	fmt.Println(record)
}
