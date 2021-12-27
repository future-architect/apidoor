package logger

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func TestUpdateLog(t *testing.T) {

	restore := flextime.Set(time.Date(2021, time.December, 27, 17, 1, 41, 0, time.UTC))
	defer restore()

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

	want := `2021-12-27T17:01:41Z,key,path,header1,header2
2021-12-27T17:01:41Z,key,path,header1,header2
`
	want = strings.ReplaceAll(want, "\r\n", "\n")

	if out := buffer.String(); out != want {
		t.Fatalf("result:\n%s\nwant:\n%s", out, want)
	}

}
