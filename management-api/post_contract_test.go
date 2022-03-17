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
	if _, err := db.Exec("DELETE FROM contract_product_content"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM contract"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Exec("DELETE FROM contract_product_content")
		db.Exec("DELETE FROM contract")
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

	tests := []struct {
		name                  string
		req                   model.PostContractReq
		wantStatus            int
		wantResp              interface{}
		wantDBID              *int
		wantContractProductDB []contractProduct
	}{
		{
			name: "create contract properly",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
				Products: []*model.ContractProducts{
					{
						ProductName: productNames[0],
						Description: "api1",
					},
				},
			},
			wantStatus: http.StatusCreated,
			wantResp:   "Created",
			wantDBID:   &userIds[0],
			wantContractProductDB: []contractProduct{
				{
					ProductID:   productIDs[0],
					Description: "api1",
				},
			},
		},
		{
			name: "products field has multiple products",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[1],
				Products: []*model.ContractProducts{
					{
						ProductName: productNames[0],
						Description: "api1",
					},
					{
						ProductName: productNames[1],
						Description: "api2",
					},
				},
			},
			wantStatus: http.StatusCreated,
			wantResp:   "Created",
			wantDBID:   &userIds[1],
			wantContractProductDB: []contractProduct{
				{
					ProductID:   productIDs[0],
					Description: "api1",
				},
				{
					ProductID:   productIDs[1],
					Description: "api2",
				},
			},
		},
		{
			name: "user item with the requested account id does not exist",
			req: model.PostContractReq{
				UserAccountID: "not_exist",
				Products: []*model.ContractProducts{
					{
						ProductName: productNames[0],
						Description: "api1",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "account_id not_exist does not exist",
			},
		},
		{
			name: "product item with the requested product name does not exist",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
				Products: []*model.ContractProducts{
					{
						ProductName: "not_exist",
						Description: "api1",
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "product_name not_exist does not exist",
			},
		},
		{
			name: "products field is missed",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "products",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            0.0,
					},
				},
			},
		},
		{
			name: "product_name field is an empty array",
			req: model.PostContractReq{
				UserAccountID: userAccountIDs[0],
				Products:      []*model.ContractProducts{},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "products",
						ConstraintType: "length_gte",
						Message:        "input array length is 0, but it must be greater than or equal to 1",
						Gte:            "1",
						Got:            0.0,
					},
				},
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
			if tt.wantDBID == nil {
				return
			}

			rows, err := db.Queryx(`SELECT id
					       				FROM contract WHERE user_id=$1 `, tt.wantDBID)
			//rows, err := db.Queryx(`SELECT user_id FROM contract`)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}

			contractID := -1
			for rows.Next() {
				err = rows.Scan(&contractID)
				if err != nil {
					t.Errorf("scan contract id failed: %v", err)
				}
			}
			if contractID == -1 {
				t.Errorf("cannot get contract id")
				return
			}

			rows, err = db.Queryx(`SELECT product_id, description
					       				FROM contract_product_content WHERE contract_id=$1 ORDER BY product_id`, contractID)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}

			gotContractProduct := make([]contractProduct, 0)
			var cp contractProduct
			for rows.Next() {
				if err = rows.StructScan(&cp); err != nil {
					t.Errorf("cannnot scan contract product: %v", err)
					return
				}
				gotContractProduct = append(gotContractProduct, cp)
			}

			if diff := cmp.Diff(tt.wantContractProductDB, gotContractProduct); diff != "" {
				t.Errorf("contract_product_content differs: \n%s", diff)
			}

		})
	}

}

type contractProduct struct {
	ProductID   int    `db:"product_id"`
	Description string `db:"description"`
}
