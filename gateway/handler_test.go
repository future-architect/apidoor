package gateway_test

import (
	"gateway"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testdata struct {
	rescode int
	content string
	apikey  string
	field   string
	request string
	out     string
	outcode int
}

var table = []testdata{
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
	// unauthorized request (invalid key)
	{
		rescode: http.StatusOK,
		content: "application/json",
		apikey:  "apikey2",
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

func TestHandler(t *testing.T) {
	for index, tt := range table {
		message := []byte("response from API server")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.rescode)
			w.Write(message)
		}))
		defer ts.Close()

		host := ts.URL[6:]

		u := gateway.NewURITemplate(tt.field)
		v := gateway.NewURITemplate(host)
		gateway.APIData["apikey1"] = append(gateway.APIData["apikey1"], gateway.Field{
			Template: *u,
			Path:     *v,
		})

		r := httptest.NewRequest(http.MethodGet, tt.request, nil)
		r.Header.Set("Content-Type", tt.content)
		r.Header.Set("Authorization", tt.apikey)
		w := httptest.NewRecorder()
		gateway.GetHandler(w, r)

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
