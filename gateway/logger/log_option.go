package logger

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Songmu/flextime"
)

type LogOption func(record *[]string, key, path string, r *http.Request)

func WithTime() LogOption {
	return func(record *[]string, key, path string, r *http.Request) {
		now := flextime.Now()
		*record = append(*record, now.Format(time.RFC3339))
	}
}

func WithKey() LogOption {
	return func(record *[]string, key, path string, r *http.Request) {
		*record = append(*record, key)
	}
}

func WithPath() LogOption {
	return func(record *[]string, key, path string, r *http.Request) {
		*record = append(*record, path)
	}
}

func HeaderElement(name string) LogOption {
	return func(record *[]string, key, path string, r *http.Request) {
		*record = append(*record, r.Header.Get(name))
	}
}

var logOptionPattern = DefaultLogPattern()

func DefaultLogPattern() []LogOption {
	// check log pattern
	ptn := os.Getenv("LOG_PATTERN")
	if ptn == "" {
		ptn = "time,key,path" // default
	}

	// get column name and set function to write log
	var pattern []LogOption
	for _, value := range strings.Split(ptn, ",") {
		switch value {
		case "time":
			pattern = append(pattern, WithTime())
		case "key":
			pattern = append(pattern, WithKey())
		case "path":
			pattern = append(pattern, WithPath())
		default:
			pattern = append(pattern, HeaderElement(value))
		}
	}

	return pattern
}
