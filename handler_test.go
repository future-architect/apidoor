package apidoor_test

import (
	"apidoor"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testdata struct {
	rescode       int
	content       string
	authorization string
	apinum        string
	apikey        string
	out           string
	outcode       int
}

var table = []testdata{
	// valid request
	{
		rescode:       http.StatusOK,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "0",
		apikey:        "apikey1",
		out:           "response from API server",
		outcode:       http.StatusOK,
	},
	// client error
	{
		rescode:       http.StatusBadRequest,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "0",
		apikey:        "apikey1",
		out:           "response from API server",
		outcode:       http.StatusBadRequest,
	},
	// server error
	{
		rescode:       http.StatusInternalServerError,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "0",
		apikey:        "apikey1",
		out:           "response from API server",
		outcode:       http.StatusInternalServerError,
	},
	// invalid Content-Type
	{
		rescode:       http.StatusOK,
		content:       "text/html",
		authorization: "testtoken",
		apinum:        "0",
		apikey:        "apikey1",
		out:           "unexpected request content",
		outcode:       http.StatusBadRequest,
	},
	// invalid Authorization
	{
		rescode:       http.StatusOK,
		content:       "application/json",
		authorization: "",
		apinum:        "0",
		apikey:        "apikey1",
		out:           "forbidden",
		outcode:       http.StatusForbidden,
	},
	// invalid apinum
	{
		rescode:       http.StatusOK,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "hoge",
		apikey:        "apikey1",
		out:           "invalid API number",
		outcode:       http.StatusBadRequest,
	},
	// unauthorized request (invalid key)
	{
		rescode:       http.StatusOK,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "0",
		apikey:        "apikey5",
		out:           "error: invalid key",
		outcode:       http.StatusBadRequest,
	},
	// unauthorized request (invalid num)
	{
		rescode:       http.StatusOK,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "-1",
		apikey:        "apikey1",
		out:           "error: unauthorized request",
		outcode:       http.StatusBadRequest,
	},
	// invalid URL
	{
		rescode:       http.StatusOK,
		content:       "application/json",
		authorization: "testtoken",
		apinum:        "1",
		apikey:        "apikey2",
		out:           "Get \"hoge\": unsupported protocol scheme \"\"",
		outcode:       http.StatusInternalServerError,
	},
}

func TestHandler(t *testing.T) {
	for index, tt := range table {
		message := []byte("response from API server")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.rescode)
			w.Write(message)
		}))
		defer ts.Close()

		// change destination of API temporarily
		apidoor.Urldata.Url = []string{
			ts.URL,
			"hoge",
		}

		r := httptest.NewRequest(http.MethodGet, "/hoge?apikey="+tt.apikey+"&num="+tt.apinum, nil)
		r.Header.Set("Content-Type", tt.content)
		r.Header.Set("Authorization", tt.authorization)
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
