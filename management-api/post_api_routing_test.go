package managementapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/go-redis/redis/v8"
)

func TestPostAPIRouting(t *testing.T) {
	targetKey := "APIKEY"

	//setup DB
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	t.Cleanup(func() {
		err := redisClient.Del(ctx, targetKey).Err()
		if err != nil {
			t.Fatalf("delete redis fields error: %v", err)
		}
	})

	tests := []struct {
		name          string
		apiKey        string
		path          string
		forwardURL    string
		checkHgetArgs []string
		checkHgetResp string
		httpStatus    int
		wantResp      interface{}
	}{
		{
			name:          "[正常系] 正しくルーティングを登録できる",
			apiKey:        targetKey,
			path:          "test",
			forwardURL:    "http://localhost/test",
			checkHgetArgs: []string{targetKey, "test"},
			checkHgetResp: "http://localhost/test",
			httpStatus:    http.StatusCreated,
			wantResp:      "Created",
		},
		{
			name:          "[異常系] 空文字列のパラメータがある場合、ルーティングを登録しない",
			apiKey:        "",
			path:          "test2",
			forwardURL:    "http://localhost/test",
			checkHgetArgs: []string{targetKey, "test2"},
			checkHgetResp: "",
			httpStatus:    http.StatusBadRequest,
			wantResp: managementapi.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &managementapi.ValidationErrors{
					{
						Field:          "api_key",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
		},
		{
			name:          "[異常系] forward_urlがURL schemeを満たしていないとき、ルーティングを登録しない",
			apiKey:        targetKey,
			path:          "test3",
			forwardURL:    "wrong_url",
			checkHgetArgs: []string{targetKey, "test3"},
			checkHgetResp: "",
			httpStatus:    http.StatusBadRequest,
			wantResp: managementapi.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &managementapi.ValidationErrors{
					{
						Field:          "forward_url",
						ConstraintType: "url",
						Message:        "input value, wrong_url, does not satisfy the format, url",
						Got:            "wrong_url",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqParam := managementapi.PostAPIRoutingReq{
				ApiKey:     tt.apiKey,
				Path:       tt.path,
				ForwardURL: tt.forwardURL,
			}
			body, err := json.Marshal(reqParam)
			if err != nil {
				t.Errorf("faild creating body: %v", err)
				return
			}
			r := httptest.NewRequest(http.MethodPost, "localhost:3001/mgmt/routing", bytes.NewReader(body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			managementapi.PostAPIRouting(w, r)
			rw := w.Result()
			defer rw.Body.Close()
			if rw.StatusCode != tt.httpStatus {
				t.Errorf("responce status is wrong: want %d, got %d", tt.httpStatus, rw.StatusCode)
			}
			if len(tt.checkHgetArgs) < 2 {
				return
			}
			val := redisClient.HGet(ctx, tt.checkHgetArgs[0], tt.checkHgetArgs[1]).Val()
			if val != tt.checkHgetResp {
				t.Errorf("Hget responce is wrong: want %s, got %s", tt.checkHgetResp, val)
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
			case managementapi.BadRequestResp:
				want := tt.wantResp.(managementapi.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not unsupported")
			}
		})
	}
}
