package logger

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"testing"
)

func TestUpdateLog(t *testing.T) {

	r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Set("TEST1", "header1")
	r.Header.Set("TEST2", "header2")
	logOptionPattern = []LogOption{
		WithTime(),
		WithKey(),
		WithPath(),
		HeaderElement("TEST1"),
		HeaderElement("TEST2"),
	}

	buffer := &bytes.Buffer{}
	appender := DefaultAppender{
		Writer: buffer,
	}
	for i := 0; i < 2; i++ {
		appender.Do("key", "path", r)
	}

	// check if log is valid
	reader := csv.NewReader(buffer)
	recordNum := 0

	t.Log(buffer.String())

	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatal(err)
		}
		if line[1] != "key" {
			t.Fatalf("unexpected log %s, expected 'key'", line[1])
		} else if line[2] != "path" {
			t.Fatalf("unexpected log %s, expected 'path'", line[2])
		} else if line[3] != "header1" {
			t.Fatalf("unexpected log %s, expected 'header1'", line[3])
		} else if line[4] != "header2" {
			t.Fatalf("unexpected log %s, expected 'header2'", line[4])
		}
		recordNum++
	}

	if recordNum != 2 {
		t.Fatalf("unexpected number of log %d, expected 2", recordNum)
	}
}
