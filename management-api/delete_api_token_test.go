package managementapi_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/guregu/dynamo"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestDeleteAPIToken(t *testing.T) {
	dbType := managementapi.GetAPIDBType(t)
	if dbType != managementapi.DYNAMO {
		log.Println("this test is valid when dynamodb is used, skip")
		return
	}

	managementapi.Setup(t,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../dynamo_table/api_routing_table.json`,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../dynamo_table/access_token_table.json`,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 batch-write-item --request-items file://./testdata/delete_api_token.json`,
	)
	t.Cleanup(func() {
		managementapi.Teardown(t,
			`aws dynamodb --profile local --endpoint-url http://localhost:4566 delete-table --table access_token`,
			`aws dynamodb --profile local --endpoint-url http://localhost:4566 delete-table --table api_routing`,
		)
	})

	dbEndpoint := os.Getenv("DYNAMO_ENDPOINT")
	db := dynamo.New(session.Must(session.NewSessionWithOptions(session.Options{
		Profile:           "local",
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Endpoint: aws.String(dbEndpoint)},
	})))

	tests := []struct {
		name           string
		apiKey         string
		path           string
		wantStatusCode int
		wantResp       interface{}
		checkDBKey     string
		wantDBResp     *accessToken
	}{
		{
			name:           "delete api token properly",
			apiKey:         "key",
			path:           "test/{correct}",
			wantStatusCode: http.StatusNoContent,
			wantResp:       nil,
			checkDBKey:     "key#test/{correct}",
			wantDBResp:     nil,
		},
		{
			name:           "no token is registered",
			apiKey:         "key",
			path:           "test/no/token/registered",
			wantStatusCode: http.StatusNoContent,
			wantResp:       nil,
			checkDBKey:     "key#test/no/token/registered",
			wantDBResp:     nil,
		},
		{
			name:           "some required field is empty",
			apiKey:         "key",
			path:           "",
			wantStatusCode: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					&validator.ValidationError{
						Field:          "path",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
			checkDBKey: "key#test/insufficient/parameters",
			wantDBResp: &accessToken{
				Key: "key#test/insufficient/parameters",
				AccessTokens: []model.AccessToken{
					{
						ParamType: model.Header,
						Key:       "token",
						Value:     "token_value",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			if tt.apiKey != "" {
				form.Add("api_key", tt.apiKey)
			}
			if tt.path != "" {
				form.Add("path", tt.path)
			}
			r := httptest.NewRequest(http.MethodGet, "localhost:3000/api/search?"+form.Encode(), nil)
			w := httptest.NewRecorder()
			managementapi.DeleteAPIToken(w, r)
			rw := w.Result()
			defer rw.Body.Close()
			if rw.StatusCode != tt.wantStatusCode {
				t.Errorf("responce status is wrong: want %d, got %d", tt.wantStatusCode, rw.StatusCode)
			}

			resp, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatal(err)
			}

			switch tt.wantResp.(type) {
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			case nil:
				if len(resp) > 0 {
					t.Errorf("expected response body is empty, got %+v", resp)
				}
			default:
				t.Errorf("type of wantResp is not unsupported")
			}

			// db checkWithContext
			accessToken, err := getAccessToken(db, tt.checkDBKey)
			if err != nil {
				t.Errorf("failed to get access token: %v", err)
				return
			}
			if diff := cmp.Diff(tt.wantDBResp, accessToken); diff != "" {
				t.Errorf("get access token result differs:\n %v", err)
			}
		})
	}
}
