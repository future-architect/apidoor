package gateway

import (
	"net/http"
	"time"
)

type LogElement interface {
	append(record []string, key, path string, r *http.Request) []string
}

type TimeElement struct {
}

func (element TimeElement) append(record []string, key, path string, r *http.Request) []string {
	record = append(record, time.Now().Format(time.RFC3339))
	return record
}

type APIKeyElement struct {
}

func (element APIKeyElement) append(record []string, key, path string, r *http.Request) []string {
	record = append(record, key)
	return record
}

type PathElement struct {
}

func (element PathElement) append(record []string, key, path string, r *http.Request) []string {
	record = append(record, path)
	return record
}

type HeaderElement struct {
	name string
}

func (element HeaderElement) append(record []string, key, path string, r *http.Request) []string {
	record = append(record, r.Header.Get(element.name))
	return record
}
