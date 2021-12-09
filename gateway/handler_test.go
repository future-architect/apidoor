package gateway

import (
	"context"
	"github.com/future-architect/apidoor/gateway/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var dbHost, templatePath string

type dbMock struct{}

func (dm dbMock) GetFields(_ context.Context, key string) (Fields, error) {
	if key == "apikeyNotExist" {
		return nil, ErrUnauthorizedRequest
	}
	return Fields{
		{
			Template: *NewURITemplate(templatePath),
			Path:     *NewURITemplate(dbHost),
			Num:      5,
			Max:      10,
		},
	}, nil
}

var methods = []string{
	http.MethodGet,
	http.MethodDelete,
	http.MethodPost,
	http.MethodPut,
}

func TestHandle(t *testing.T) {

	type tests struct {
		resCode int
		content string
		apikey  string
		field   string
		request string
		out     string
		outCode int
	}

	var cases = []tests{
		// valid request using parameter
		{
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "response from API server",
			outCode: http.StatusOK,
		},
		// valid request not using parameter
		{
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test",
			request: "/test",
			out:     "response from API server",
			outCode: http.StatusOK,
		},
		// client error
		{
			resCode: http.StatusBadRequest,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "client error",
			outCode: http.StatusBadRequest,
		},
		// server error
		{
			resCode: http.StatusInternalServerError,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "server error",
			outCode: http.StatusInternalServerError,
		},
		// invalid Content-Type
		{
			resCode: http.StatusOK,
			content: "text/html",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "unexpected request content",
			outCode: http.StatusBadRequest,
		},
		// no authorization header
		{
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "no authorization",
			outCode: http.StatusBadRequest,
		},
		// unauthorized request (invalid key)
		{
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikeyNotExist",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "invalid key or path",
			outCode: http.StatusNotFound,
		},
		// unauthorized request (invalid URL)
		{
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/t/hoge",
			out:     "invalid key or path",
			outCode: http.StatusNotFound,
		},
	}

	handler := DefaultHandler{
		Appender: logger.DefaultAppender{
			Writer: os.Stdout,
		},
		DataSource: dbMock{},
	}

	for _, method := range methods {
		for index, tt := range cases {
			// http server for test
			message := []byte("response from API server")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.resCode)
				w.Write(message)
			}))

			// set routing data
			host := ts.URL[6:]
			dbHost = host
			templatePath = tt.field

			// send request to test server
			r := httptest.NewRequest(method, tt.request, nil)
			r.Header.Set("Content-Type", tt.content)
			if tt.apikey != "" {
				r.Header.Set("Authorization", tt.apikey)
			}
			w := httptest.NewRecorder()
			handler.Handle(w, r)

			// check response
			rw := w.Result()

			if rw.StatusCode != tt.outCode {
				t.Fatalf("method %s, case %d: unexpected status code %d, expected %d", method, index, rw.StatusCode, tt.outCode)
			}

			b, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatalf("method %s, case %d: unexpected body type", method, index)
			}

			trimmed := strings.TrimSpace(string(b))
			if trimmed != tt.out {
				t.Fatalf("method %s, case %d: unexpected response: %s, expected: %s", method, index, trimmed, tt.out)
			}

			// loopの中なのでdeferは使えない
			ts.Close()
			rw.Body.Close()
		}
	}
}
