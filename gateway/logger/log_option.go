package logger

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)

type LogOption func(record *[]string, logItem *LogItem, r *http.Request)

func WithTime() LogOption {
	return func(record *[]string, logItem *LogItem, r *http.Request) {
		*record = append(*record, logItem.TimeStamp)
	}
}

func WithKey() LogOption {
	return func(record *[]string, logItem *LogItem, r *http.Request) {
		*record = append(*record, logItem.Key)
	}
}

func WithPath() LogOption {
	return func(record *[]string, logItem *LogItem, r *http.Request) {
		*record = append(*record, logItem.Path)
	}
}

func WithResponseStatus() LogOption {
	return func(record *[]string, logItem *LogItem, r *http.Request) {
		*record = append(*record, strconv.Itoa(logItem.StatusCode))
	}
}

func WithBillingStatus() LogOption {
	return func(record *[]string, logItem *LogItem, r *http.Request) {
		*record = append(*record, logItem.BillingStatus.String())
	}
}

func HeaderElement(name string) LogOption {
	return func(record *[]string, logItem *LogItem, r *http.Request) {
		*record = append(*record, r.Header.Get(name))
	}
}

var LogOptionPattern = DefaultLogPattern()

func DefaultLogPattern() []LogOption {
	// check log pattern
	ptn := os.Getenv("LOG_PATTERN")
	if ptn == "" {
		ptn = "time,key,path,response_status,billing_status" // default
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
		case "response_status":
			pattern = append(pattern, WithResponseStatus())
		case "billing_status":
			pattern = append(pattern, WithBillingStatus())
		default:
			pattern = append(pattern, HeaderElement(value))
		}
	}

	return pattern
}
