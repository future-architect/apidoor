package gateway_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/future-architect/apidoor/gateway"
)

var dbHost, templatePath string

type dbMock struct{}

func (dm dbMock) GetFields(ctx context.Context, key string) (gateway.Fields, error) {
	if key == "apikeyNotExist" {
		return nil, gateway.ErrUnauthorizedRequest
	}
	return gateway.Fields{
		{
			Template: *gateway.NewURITemplate(templatePath),
			Path:     *gateway.NewURITemplate(dbHost),
			Num:      5,
			Max:      10,
		},
	}, nil
}

type handlerTest struct {
	rescode int
	content string
	apikey  string
	field   string
	request string
	out     string
	outcode int
}

var handlerTestData = []handlerTest{
	// valid request using parameter
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikey1",
		field:   "/test/{test}",
		request: "/test/hoge",
		out:     "response from API server",
		outcode: http.StatusOK,
	},
	// valid request not using parameter
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikey1",
		field:   "/test",
		request: "/test",
		out:     "response from API server",
		outcode: http.StatusOK,
	},
	// client error
	{
		rescode: http.StatusBadRequest,
		content: "application/json",
		apikey:  "apikey1",
		field:   "/test/{test}",
		request: "/test/hoge",
		out:     "client error",
		outcode: http.StatusBadRequest,
	},
	// server error
	{
		rescode: http.StatusInternalServerError,
		content: "application/json",
		apikey:  "apikey1",
		field:   "/test/{test}",
		request: "/test/hoge",
		out:     "server error",
		outcode: http.StatusInternalServerError,
	},
	// invalid Content-Type
	{
		rescode: http.StatusOK,
		content: "text/html",
		apikey:  "apikey1",
		field:   "/test/{test}",
		request: "/test/hoge",
		out:     "unexpected request content",
		outcode: http.StatusBadRequest,
	},
	// no authorization header
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "",
		field:   "/test/{test}",
		request: "/test/hoge",
		out:     "no authorization",
		outcode: http.StatusBadRequest,
	},
	// unauthorized request (invalid key)
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikeyNotExist",
		field:   "/test/{test}",
		request: "/test/hoge",
		out:     "invalid key or path",
		outcode: http.StatusNotFound,
	},
	// unauthorized request (invalid URL)
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikey1",
		field:   "/test/{test}",
		request: "/t/hoge",
		out:     "invalid key or path",
		outcode: http.StatusNotFound,
	},
}

var methods = []string{
	http.MethodGet,
	http.MethodDelete,
	http.MethodPost,
	http.MethodPut,
}

func TestHandler(t *testing.T) {
	mock := dbMock{}
	gateway.DBDriver = mock
	for _, method := range methods {
		for index, tt := range handlerTestData {
			// http server for test
			message := []byte("response from API server")
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.rescode)
				w.Write(message)
			}))
			defer ts.Close()

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
			gateway.Handler(w, r)

			// check response
			rw := w.Result()
			defer rw.Body.Close()

			if rw.StatusCode != tt.outcode {
				t.Fatalf("method %s, case %d: unexpected status code %d, expected %d", method, index, rw.StatusCode, tt.outcode)
			}

			b, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatalf("method %s, case %d: unexpected body type", method, index)
			}

			trimmed := strings.TrimSpace(string(b))
			if trimmed != tt.out {
				t.Fatalf("method %s, case %d: unexpected response: %s, expected: %s", method, index, trimmed, tt.out)
			}
		}
	}
}
