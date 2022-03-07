package managementapi_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"os"
	"testing"
)

func TestPostAPIToken(t *testing.T) {
	dbType := managementapi.GetAPIDBType(t)
	if dbType != managementapi.DYNAMO {
		log.Println("this test is valid when dynamodb is used, skip")
		return
	}

	managementapi.Setup(t,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../dynamo_table/api_routing_table.json`,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../dynamo_table/access_token_table.json`,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 batch-write-item --request-items file://./testdata/post_api_token.json`,
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
		req            model.PostAPITokenReq
		wantStatusCode int
		wantResp       interface{}
		checkDBKey     string
		wantDBResp     *accessToken
	}{
		{
			name: "post api token properly",
			req: model.PostAPITokenReq{
				APIKey: "key",
				Path:   "test/{correct}",
				AccessTokens: []model.AccessToken{
					{
						ParamType: model.Header,
						Key:       "token",
						Value:     "token_value",
					},
				},
			},
			wantStatusCode: http.StatusCreated,
			wantResp:       "Created",
			checkDBKey:     "key#test/{correct}",
			wantDBResp: &accessToken{
				Key: "key#test/{correct}",
				AccessTokens: []model.AccessToken{
					{
						ParamType: model.Header,
						Key:       "token",
						Value:     "token_value",
					},
				},
			},
		},
		{
			name: "do not post api token when a pair of api_key and path does not exist",
			req: model.PostAPITokenReq{
				APIKey: "key",
				Path:   "test/not/exist",
				AccessTokens: []model.AccessToken{
					{
						ParamType: model.Header,
						Key:       "token",
						Value:     "token_value",
					},
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantResp:       validator.BadRequestResp{Message: "api_key or path is wrong"},
			checkDBKey:     "key#test/not/exist",
			wantDBResp:     nil,
		},
		{
			name: "some required field is empty",
			req: model.PostAPITokenReq{
				APIKey: "key",
				Path:   "test/{correct}",
				AccessTokens: []model.AccessToken{
					{
						ParamType: "wrong_type",
						Key:       "token",
						Value:     "token_value",
					},
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					&validator.ValidationError{
						Field:          "tokens[0].param_type",
						ConstraintType: "enum",
						Message:        "input value is wrong_type, but it must be one of the following values: [header query body_from_encoded]",
						Enum:           []string{"header", "query", "body_from_encoded"},
						Got:            "wrong_type",
					},
				},
			},
			checkDBKey: "key#test/not/exist",
			wantDBResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("failed to create request body: %v", err)
				return
			}
			r := httptest.NewRequest(http.MethodPost, "localhost:3001/mgmt/api/token", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			managementapi.PostAPIToken(w, r)
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
			case string:
				if tt.wantResp != string(resp) {
					t.Errorf("response body is not %s, got %s", tt.wantResp, string(resp))
				}
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
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

func getAccessToken(db *dynamo.DB, key string) (*accessToken, error) {
	var resp *accessToken
	accessTokenTable := os.Getenv("DYNAMO_TABLE_ACCESS_TOKEN")
	err := db.Table(accessTokenTable).
		Get("key", key).
		One(&resp)

	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return resp, nil
}

type accessToken struct {
	Key          string              `dynamo:"key"`
	AccessTokens []model.AccessToken `dynamo:"tokens"`
}
