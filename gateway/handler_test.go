package gateway

import (
	"context"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/future-architect/apidoor/gateway/model"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var dbHost, templatePath string

type dbMock struct{}

func (dm dbMock) GetFields(_ context.Context, key string) (model.Fields, error) {
	if key == "apikeyNotExist" {
		return nil, model.ErrUnauthorizedRequest
	}
	return model.Fields{
		{
			ForwardSchema: "http",
			Template:      model.NewURITemplate(templatePath),
			Path:          model.NewURITemplate(dbHost),
			Num:           5,
			Max:           10,
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

	cases := []struct {
		name    string
		resCode int
		content string
		apikey  string
		field   string
		request string
		out     string
		outCode int
	}{
		{
			name:    "valid request using parameter",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "response from API server",
			outCode: http.StatusOK,
		},
		{
			name:    "valid request not using parameter",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test",
			request: "/test",
			out:     "response from API server",
			outCode: http.StatusOK,
		},
		{
			name:    "client error",
			resCode: http.StatusBadRequest,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "client error",
			outCode: http.StatusBadRequest,
		},
		{
			name:    "server error",
			resCode: http.StatusInternalServerError,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "server error",
			outCode: http.StatusInternalServerError,
		},
		{
			name:    "invalid Content-Type",
			resCode: http.StatusOK,
			content: "text/html",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "unexpected request content",
			outCode: http.StatusBadRequest,
		},
		{
			name:    "no authorization header",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "no authorization request header",
			outCode: http.StatusBadRequest,
		},
		{
			name:    "unauthorized request (invalid key)",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikeyNotExist",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "invalid key or path",
			outCode: http.StatusNotFound,
		},
		{
			name:    "unauthorized request (invalid URL)",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/t/hoge",
			out:     "invalid key or path",
			outCode: http.StatusNotFound,
		},
	}

	h := DefaultHandler{
		Appender: logger.DefaultAppender{
			Writer: os.Stdout,
		},
		DataSource: dbMock{},
	}

	for _, method := range methods {
		for _, tt := range cases {

			// http server for test
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.resCode)
				w.Write([]byte("response from API server"))
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
			h.Handle(w, r)

			// check response
			rw := w.Result()

			b, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatalf("method:%s, case:%s: unexpected body type", method, tt.name)
			}

			if rw.StatusCode != tt.outCode {
				t.Fatalf("method:%s, case:%s: unexpected status code %d, expected %d, body:%s", method, tt.name, rw.StatusCode, tt.outCode, b)
			}

			trimmed := strings.TrimSpace(string(b))
			if trimmed != tt.out {
				t.Fatalf("method:%s, case:%s: unexpected response: %s, expected: %s", method, tt.name, trimmed, tt.out)
			}

			// loopの中なのでdeferは使えない
			ts.Close()
			rw.Body.Close()
		}
	}
}
