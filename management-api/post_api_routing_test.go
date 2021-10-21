package managementapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	}{
		{
			name:          "[正常系] 正しくルーティングを登録できる",
			apiKey:        targetKey,
			path:          "test",
			forwardURL:    "localhost:3000/test",
			checkHgetArgs: []string{targetKey, "test"},
			checkHgetResp: "localhost:3000/test",
			httpStatus:    http.StatusCreated,
		},
		{
			name:          "[異常系] 空文字列のパラメータがある場合、ルーティングを登録しない",
			apiKey:        targetKey,
			path:          "test2",
			forwardURL:    "",
			checkHgetArgs: []string{targetKey, "test2"},
			checkHgetResp: "",
			httpStatus:    http.StatusBadRequest,
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
			r := httptest.NewRequest(http.MethodPost, "localhost:3001/mgmt/api/", bytes.NewReader(body))
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
		})
	}
}
