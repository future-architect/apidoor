package managementapi_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostAPIKey(t *testing.T) {
	if _, err := db.Exec("DELETE FROM apikey"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Exec("DELETE FROM apikey")
		db.Exec("DELETE FROM apiuser")
	}()

	// DB set up
	userAccountIDs := []string{"user1", "user2"}
	userIds := make([]int, len(userAccountIDs))
	for i, name := range userAccountIDs {
		stmt, err := db.Preparex(
			`INSERT INTO apiuser(account_id, email_address, login_password_hash, name, created_at, updated_at)
			VALUES ($1, 'a', 'password', 'a', current_timestamp, current_timestamp) RETURNING  id`)
		if err != nil {
			t.Error(err)
			return
		}
		var id int
		stmt.QueryRowx(name).Scan(&id)
		userIds[i] = id
	}

	tests := []struct {
		name        string
		req         model.PostAPIKeyReq
		wantStatus  int
		wantResp    interface{}
		needCheckDB bool
	}{
		{
			name: "create api key properly",
			req: model.PostAPIKeyReq{
				UserAccountID: userAccountIDs[0],
			},
			wantStatus: http.StatusCreated,
			wantResp: model.PostAPIKeyResp{
				UserAccountID: userAccountIDs[0],
			},
		},
		{
			name: "account id in request does not exist",
			req: model.PostAPIKeyReq{
				UserAccountID: "does_not_exist",
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "account_id does_not_exist does not exist",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("create request body error: %v", err)
				return
			}
			body := bytes.NewReader(bodyBytes)

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/keys", body)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()
			managementapi.PostAPIKey(w, r)

			rw := w.Result()

			resp, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Errorf("read response body error: %v", err)
				return
			}

			if rw.StatusCode != tt.wantStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantStatus)
			}

			switch tt.wantResp.(type) {
			case model.PostAPIKeyResp:
				want := tt.wantResp.(model.PostAPIKeyResp)
				gotBody := new(model.PostAPIKeyResp)
				if err := json.Unmarshal(resp, gotBody); err != nil {
					t.Errorf("parsing body as model.PostAPIKeyResp failed: %v\ngot: %v", err, string(resp))
					return
				}
				testPostAPIKeyResp(t, &want, gotBody)
				testPostAPIKeyDB(t, gotBody)
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not supported")
			}
		})
	}
}

func testPostAPIKeyResp(t *testing.T, want, got *model.PostAPIKeyResp) {
	t.Helper()

	if diff := cmp.Diff(got, want, cmpopts.IgnoreFields(*got, "AccessKey", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Errorf("PostAPIKeyResp response differs:\n%v", diff)
	}

	key, err := hex.DecodeString(got.AccessKey)
	if err != nil {
		t.Errorf("parsing access key as hex failed, key %s: %v", got.AccessKey, err)
	}
	wantKeyLength := 16
	if len(key) != wantKeyLength {
		t.Errorf("length of the access key is not %d, got %d", wantKeyLength, len(key))
	}
}

func testPostAPIKeyDB(t *testing.T, got *model.PostAPIKeyResp) {
	t.Helper()
	rows, err := db.Queryx(`SELECT access_key
       				FROM apikey WHERE access_key=$1`, got.AccessKey)
	if err != nil {
		t.Errorf("db get api info error: %v", err)
		return
	}
	cnt := 0
	for rows.Next() {
		cnt++
	}
	if cnt != 1 {
		t.Errorf("the number of db records is wrong, want 1, got %d", cnt)
	}
}
