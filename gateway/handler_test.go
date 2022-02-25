package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/gateway/logger"
	"github.com/future-architect/apidoor/gateway/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func (dm dbMock) GetAccessTokens(ctx context.Context, apikey, templatePath string) (*model.AccessTokens, error) {
	key := fmt.Sprintf("%s#%s", apikey, templatePath)
	log.Println(key)
	switch key {
	case "key#testheader":
		return &model.AccessTokens{
			Tokens: []model.AccessToken{
				{
					ParamType: model.Header,
					Key:       "Token",
					Value:     "token_value",
				},
			},
		}, nil
	case "key#testquery":
		return &model.AccessTokens{
			Tokens: []model.AccessToken{
				{
					ParamType: model.Query,
					Key:       "token",
					Value:     "token_value",
				},
			},
		}, nil
	case "key#testform":
		return &model.AccessTokens{
			Tokens: []model.AccessToken{
				{
					ParamType: model.BodyFormEncoded,
					Key:       "token",
					Value:     "token_value",
				},
			},
		}, nil
	case "key#test/unsupport":
		return &model.AccessTokens{
			Tokens: []model.AccessToken{
				{
					ParamType: "unsupported",
					Key:       "token",
					Value:     "token_value",
				},
			},
		}, nil
	case "key#test/wrongtype":
		return &model.AccessTokens{
			Tokens: []model.AccessToken{
				{
					ParamType: model.BodyFormEncoded,
					Key:       "token",
					Value:     "token_value",
				},
			},
		}, nil
	case "key#test/multiple":
		return &model.AccessTokens{
			Tokens: []model.AccessToken{
				{
					ParamType: model.Query,
					Key:       "token2",
					Value:     "token_value2",
				},
				{
					ParamType: model.Header,
					Key:       "token",
					Value:     "token_value",
				},
			},
		}, nil
	}

	return nil, nil
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
			out:     "response from API server",
			outCode: http.StatusBadRequest,
		},
		{
			name:    "server error",
			resCode: http.StatusInternalServerError,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "response from API server",
			outCode: http.StatusInternalServerError,
		},
		{
			name:    "no authorization header",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "gateway error: no authorization request header",
			outCode: http.StatusBadRequest,
		},
		{
			name:    "unauthorized request (invalid key)",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikeyNotExist",
			field:   "/test/{test}",
			request: "/test/hoge",
			out:     "gateway error: invalid key or path",
			outCode: http.StatusNotFound,
		},
		{
			name:    "unauthorized request (invalid URL)",
			resCode: http.StatusOK,
			content: "application/json",
			apikey:  "apikey1",
			field:   "/test/{test}",
			request: "/t/hoge",
			out:     "gateway error: invalid key or path",
			outCode: http.StatusNotFound,
		},
	}

	h := DefaultHandler{
		Appender: &logger.DefaultAppender{
			Writer: os.Stdout,
		},
		DataSource: dbMock{},
	}

	for _, method := range methods {
		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
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
					r.Header.Set("X-Apidoor-Authorization", tt.apikey)
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
			})
		}
	}
}

func TestSetStoredTokens(t *testing.T) {
	apikey := "key"
	tests := []struct {
		name          string
		templatePath  string
		requestMethod string
		requestURL    string      //including query param
		requestHeader http.Header //header except X-Apidoor-Authorization
		requestBody   io.Reader
		wantURL       string
		wantHeader    http.Header //header except X-Apidoor-Authorization
		wantBody      interface{}
		wantErr       error
	}{
		{
			name:          "append header parameter properly",
			templatePath:  "testheader",
			requestMethod: "GET",
			requestURL:    "http://example.com/testheader",
			wantURL:       "http://example.com/testheader",
			wantHeader:    http.Header{"Token": []string{"token_value"}},
		},
		{
			name:          "do not overwrite existing header",
			templatePath:  "testheader",
			requestMethod: "GET",
			requestHeader: http.Header{"Token": []string{"token_value_original"}},
			requestURL:    "http://example.com/testheader",
			wantURL:       "http://example.com/testheader",
			wantHeader:    http.Header{"Token": []string{"token_value_original"}},
		},
		{
			name:          "append query parameter properly",
			templatePath:  "testquery",
			requestMethod: "GET",
			requestURL:    "http://example.com/testquery",
			wantURL:       "http://example.com/testquery?token=token_value",
		},
		{
			name:          "do not overwrite existing query",
			templatePath:  "testquery",
			requestMethod: "GET",
			requestURL:    "http://example.com/testquery?token=token_value_original",
			wantURL:       "http://example.com/testquery?token=token_value_original",
		},
		{
			name:          "append query parameter properly",
			templatePath:  "testquery",
			requestMethod: "GET",
			requestURL:    "http://example.com/testquery",
			wantURL:       "http://example.com/testquery?token=token_value",
		},
		{
			name:          "do not overwrite existing form value",
			templatePath:  "testform",
			requestMethod: "POST",
			requestHeader: http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			requestBody:   createFormURLEncodedBody(map[string]string{"token": "token_value_original"}),
			requestURL:    "http://example.com/testform",
			wantURL:       "http://example.com/testform",
			wantHeader:    http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			wantBody: url.Values{
				"token": {"token_value_original"},
			},
		},
		{
			name:          "append multiple tokens properly",
			templatePath:  "test/multiple",
			requestMethod: "GET",
			requestURL:    "http://example.com/test/multiple",
			wantURL:       "http://example.com/test/multiple?token2=token_value2",
			wantHeader:    http.Header{"Token": []string{"token_value"}},
		},
		{
			name:          "no token is stored",
			templatePath:  "test/notoken",
			requestMethod: "GET",
			requestURL:    "http://example.com/test/notoken",
			wantURL:       "http://example.com/test/notoken",
		},
		{
			name:          "append header parameter properly",
			templatePath:  "test/unsupport",
			requestMethod: "GET",
			requestURL:    "http://example.com/test/unsupport",
			wantURL:       "http://example.com/test/unsupport",
			wantErr:       errors.New("unsupported param type: unsupported"),
		},
		{
			name:          "append header parameter properly",
			templatePath:  "test/wrongtype",
			requestMethod: "GET",
			requestHeader: http.Header{"Content-Type": []string{"application/json"}},
			requestURL:    "http://example.com/test/wrongtype",
			wantURL:       "http://example.com/test/wrongtype",
			wantErr:       errors.New("content-Type header is not application/x-www-form-urlencoded, got application/json"),
		},
	}
	h := DefaultHandler{
		Appender: &logger.DefaultAppender{
			Writer: os.Stdout,
		},
		DataSource: dbMock{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.requestMethod, tt.requestURL, tt.requestBody)
			if err != nil {
				t.Errorf("creating request failed: %v", err)
				return
			}
			if tt.requestHeader != nil {
				req.Header = tt.requestHeader
			}
			req.Header.Add("X-Apidoor-Authorization", apikey)

			err = setStoredTokens(context.Background(), tt.templatePath, req, h.DataSource)
			if err != nil {
				if errors.Is(err, tt.wantErr) {
					t.Errorf("returned error differs: want %v, got %v", tt.wantErr, err)
				}
				return
			}
			if tt.wantErr != nil {
				t.Errorf("exptected error is %v, got nil", tt.wantErr)
				return
			}

			// validate the updated request
			if req.URL.String() != tt.wantURL {
				t.Errorf("updated URL differs: want %v, got %v", tt.wantURL, req.URL)
			}

			if tt.wantHeader == nil {
				req.Header.Del("X-Apidoor-Authorization")
				if len(req.Header) > 0 {
					t.Errorf("updated header is not nil, got %v", req.Header)
				}
			} else if diff := cmp.Diff(tt.wantHeader, req.Header, cmpopts.IgnoreMapEntries(func(k string, v []string) bool {
				return k == "X-Apidoor-Authorization"
			})); diff != "" {
				t.Errorf("updated header differs: \n%v", diff)
			}

			switch req.Header.Get("Content-Type") {
			case "application/x-www-form-urlencoded":
				err = req.ParseForm()
				if err != nil {
					t.Errorf("cannot parse body as form: %v", err)
					break
				}
				form := req.PostForm
				wantForm, _ := tt.wantBody.(url.Values)
				if diff := cmp.Diff(wantForm, form); diff != "" {
					t.Errorf("form in body differs: \n%v", diff)
				}
			}

		})
	}

}

func createFormURLEncodedBody(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}
