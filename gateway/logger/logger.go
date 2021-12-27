package logger

import (
	"encoding/csv"
	"io"
	"net/http"
	"strings"
)

type Appender interface {
	Do(key, path string, r *http.Request) error
}

type DefaultAppender struct {
	Writer io.Writer
}

func (a DefaultAppender) Do(key, path string, r *http.Request) error {
	record := make([]string, 0, len(logOptionPattern))
	for _, logOption := range logOptionPattern {
		logOption(&record, key, path, r)
	}

	// デフォルトはカンマ区切り
	_, err := a.Writer.Write([]byte(strings.Join(record, ",") + "\n"))

	return err
}

type CSVAppender struct {
	Writer *csv.Writer
}

func (a CSVAppender) Do(key, path string, r *http.Request) error {
	record := make([]string, 0, len(logOptionPattern))
	for _, logOption := range logOptionPattern {
		logOption(&record, key, path, r)
	}
	return a.Writer.Write(record)
}
