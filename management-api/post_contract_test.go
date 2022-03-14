package managementapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostContract(t *testing.T) {
	if _, err := db.Exec("TRUNCATE contract"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Exec("TRUNCATE contract")
		db.Exec("DELETE FROM apiuser")
		db.Exec("DELETE FROM product")
	}()

	// DB setup
	productNames := []string{"product1", "product2"}
	productIDs := make([]int, len(productNames))
	for i, name := range productNames {
		stmt, err := db.Preparex(
			`INSERT INTO product(name, source, description, thumbnail, display_name, created_at, updated_at)
			VALUES ($1, 'a', 'a', 'a', 'a', current_timestamp, current_timestamp) RETURNING id`)
		if err != nil {
			t.Error(err)
			return
		}
		var id int
		stmt.QueryRowx(name).Scan(&id)
		productIDs[i] = id
	}

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

	type checkDBValues struct {
		userID    int
		productID int
	}

	tests := []struct {
		name          string
		req           model.PostContractReq
		wantStatus    int
		wantResp      interface{}
		checkDBValues *checkDBValues
		wantDBResp    *model.Contract
	}{
		{
			name: "create contract properly",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
				ProductName:   productNames[0],
			},
			wantStatus: http.StatusCreated,
			wantResp:   "Created",
			checkDBValues: &checkDBValues{
				userID:    userIds[0],
				productID: productIDs[0],
			},
			wantDBResp: &model.Contract{
				UserID:    userIds[0],
				ProductID: productIDs[0],
			},
		},
		{
			name: "user item with the requested account id does not exist",
			req: model.PostContractReq{
				UserAccountID: "not_exist",
				ProductName:   productNames[0],
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "account_id not_exist does not exist",
			},
			checkDBValues: nil,
			wantDBResp:    nil,
		},
		{
			name: "product item with the requested product name does not exist",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
				ProductName:   "not_exist",
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "product_name not_exist does not exist",
			},
			checkDBValues: nil,
			wantDBResp:    nil,
		},
		{
			name: "product_name field is missed",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "product_name",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
			checkDBValues: nil,
			wantDBResp:    nil,
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

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/contracts", body)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()
			managementapi.PostContract(w, r)

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
			case string:
				if tt.wantResp != string(resp) {
					t.Errorf("wrong reponse body: got %s, want %s", resp, tt.wantResp)
				}
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not supported")
			}

			// db check
			if tt.checkDBValues == nil {
				return
			}

			rows, err := db.Queryx(`SELECT user_id, product_id
       				FROM contract WHERE user_id=$1 and product_id=$2`, tt.checkDBValues.userID, tt.checkDBValues.productID)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}

			var contract model.Contract
			for rows.Next() {
				if err := rows.StructScan(&contract); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}
			}
			if diff := cmp.Diff(*tt.wantDBResp, contract); diff != "" {
				t.Errorf("db get contract responce differs:\n %v", diff)
			}
		})
	}

}
