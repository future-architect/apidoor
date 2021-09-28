package gateway

import (
	"net/http"
	"time"
)

type LogOption func(record *[]string, key, path string, r *http.Request)

func WithTime() LogOption {
	return func(record *[]string, key, path string, r *http.Request) {
		*record = append(*record, time.Now().Format(time.RFC3339))
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
