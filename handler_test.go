package apidoor_test

import (
	"apidoor"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
)

type testdata struct {
	rescode int
	content string
	apikey  string
	out     string
	outcode int
}

var table = []testdata{
	// valid request
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikey1",
		out:     "response from API server",
		outcode: http.StatusOK,
	},
	// client error
	{
		rescode: http.StatusBadRequest,
		content: "application/json",
		apikey:  "apikey1",
		out:     "response from API server",
		outcode: http.StatusBadRequest,
	},
	// server error
	{
		rescode: http.StatusInternalServerError,
		content: "application/json",
		apikey:  "apikey1",
		out:     "response from API server",
		outcode: http.StatusInternalServerError,
	},
	// invalid Content-Type
	{
		rescode: http.StatusOK,
		content: "text/html",
		apikey:  "apikey1",
		out:     "unexpected request content",
		outcode: http.StatusBadRequest,
	},
	// unauthorized request (invalid key)
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikey2",
		out:     "error: unauthorized request",
		outcode: http.StatusBadRequest,
	},
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func TestHandler(t *testing.T) {
	for index, tt := range table {
		message := []byte("response from API server")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.rescode)
			w.Write(message)
		}))
		defer ts.Close()

		host := ts.URL[6:]

		// change destination of API temporarily
		rdb.FlushAll(ctx)
		rdb.SAdd(ctx, "apikey1", host)

		r := httptest.NewRequest(http.MethodGet, host, nil)
		r.Header.Set("Content-Type", tt.content)
		r.Header.Set("Authorization", tt.apikey)
		w := httptest.NewRecorder()
		apidoor.Handler(w, r)

		rw := w.Result()
		defer rw.Body.Close()

		if rw.StatusCode != tt.outcode {
			t.Fatalf("case %d: unexpected status code %d, expected %d", index, rw.StatusCode, tt.outcode)
		}

		b, err := io.ReadAll(rw.Body)
		if err != nil {
			t.Fatalf("case %d: unexpected body type", index)
		}

		trimmed := strings.TrimSpace(string(b))
		if trimmed != tt.out {
			t.Fatalf("case %d: unexpected response: %s, expected: %s", index, trimmed, tt.out)
		}

	}
}
