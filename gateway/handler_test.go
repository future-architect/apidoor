package gateway_test

import (
	"context"
	"fmt"
	"gateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
)

var (
	// redis
	rdb *redis.Client
	// dynamo
	ddb                *dynamo.DB
	apiForwardingTable string
)

var ctx = context.Background()

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

type handlerData struct {
	handler func(http.ResponseWriter, *http.Request)
	method  string
}

var handlerList = []handlerData{
	{
		handler: gateway.GetHandler,
		method:  "GET",
	},
	{
		handler: gateway.PostHandler,
		method:  "POST",
	},
	{
		handler: gateway.PutHandler,
		method:  "PUT",
	},
	{
		handler: gateway.DeleteHandler,
		method:  "DELETE",
	},
}

func setupDB(t *testing.T, dbType string) {
	if dbType == "REDIS" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
			Password: "",
			DB:       0,
		})
	} else if dbType == "DYNAMO" {
		ddb = dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Endpoint: aws.String(os.Getenv("DYNAMO_ENDPOINT"))},
		})))
		apiForwardingTable = os.Getenv("DYNAMO_TABLE_API_FORWARDING")
	} else {
		t.Fatalf("invalid db type: %s", dbType)
	}
}

func TestHandler(t *testing.T) {
	dbType := os.Getenv("DB_TYPE")
	setupDB(t, dbType)
	for _, h := range handlerList {
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
			if dbType == "REDIS" {
				rdb.HSet(ctx, "apikey1", tt.field, host)
			} else if dbType == "DYNAMO" {
				item := gateway.APIForwarding{
					APIKey:     "apikey1",
					Path:       tt.field,
					ForwardURL: host,
				}
				err := ddb.Table(apiForwardingTable).
					Put(item).Run()
				if err != nil {
					t.Errorf("put item error: %v", err)
				}
			}

			// send request to test server
			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			r.Header.Set("Content-Type", tt.content)
			if tt.apikey != "" {
				r.Header.Set("Authorization", tt.apikey)
			}
			w := httptest.NewRecorder()
			h.handler(w, r)

			// check response
			rw := w.Result()
			defer rw.Body.Close()

			if rw.StatusCode != tt.outcode {
				t.Fatalf("method %s, case %d: unexpected status code %d, expected %d", h.method, index, rw.StatusCode, tt.outcode)
			}

			b, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatalf("method %s, case %d: unexpected body type", h.method, index)
			}

			trimmed := strings.TrimSpace(string(b))
			if trimmed != tt.out {
				t.Fatalf("method %s, case %d: unexpected response: %s, expected: %s", h.method, index, trimmed, tt.out)
			}
		}
	}
}
