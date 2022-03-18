package managementapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPostAPIKeyProducts(t *testing.T) {
	if _, err := db.Exec("TRUNCATE apikey_contract_product_authorized"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM contract_product_content"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM contract"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM apikey"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Exec("TRUNCATE apikey_contract_product_authorized")
		db.Exec("DELETE FROM contract_product_content")
		db.Exec("DELETE FROM contract")
		db.Exec("DELETE FROM product")
		db.Exec("DELETE FROM apikey")
		db.Exec("DELETE FROM apiuser")
	}()

	// DB setup
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

	productNames := []string{"product1", "product2", "product3", "product4"}
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

	contractUserIDs := []int{userIds[0], userIds[1], userIds[1]}
	contractIDs := make([]int, len(contractUserIDs))
	for i, userID := range contractUserIDs {
		stmt, err := db.Preparex(
			`INSERT INTO contract( user_id, created_at, updated_at)
			VALUES ($1, current_timestamp, current_timestamp) RETURNING id`)
		if err != nil {
			t.Error(err)
			return
		}
		var id int
		stmt.QueryRowx(userID).Scan(&id)
		contractIDs[i] = id
	}

	apikeyUserIDs := []int{userIds[0], userIds[1]}
	apikeyIDs := make([]int, len(apikeyUserIDs))
	for i, userID := range apikeyUserIDs {
		key := strconv.Itoa(i)
		stmt, err := db.Preparex(
			`INSERT INTO apikey( user_id, access_key, created_at, updated_at)
			VALUES ($1, $2, current_timestamp, current_timestamp) RETURNING id`)
		if err != nil {
			t.Error(err)
			return
		}
		var id int
		stmt.QueryRowx(userID, key).Scan(&id)
		apikeyIDs[i] = id
	}

	contractProductContentContractIDs := []int{contractIDs[0], contractIDs[0], contractIDs[1], contractIDs[1], contractIDs[2], contractIDs[2]}
	contractProductContentProductIDs := []int{productIDs[0], productIDs[1], productIDs[0], productIDs[1], productIDs[2], productIDs[3]}
	contractProductContentIDs := make([]int, len(contractProductContentContractIDs))
	for i := range contractProductContentContractIDs {
		stmt, err := db.Preparex(
			`INSERT INTO contract_product_content( contract_id, product_id, created_at, updated_at)
			VALUES ($1, $2, current_timestamp, current_timestamp) RETURNING id`)
		if err != nil {
			t.Error(err)
			return
		}
		var id int
		stmt.QueryRowx(contractProductContentContractIDs[i], contractProductContentProductIDs[i]).Scan(&id)
		contractProductContentIDs[i] = id
	}

	type dbKeys struct {
		apiKeyID          int
		contractProductID []int
	}

	notExistID := -1

	tests := []struct {
		name       string
		req        model.PostAPIKeyProductsReq
		wantStatus int
		wantResp   interface{}
		wantDBKeys *dbKeys
	}{
		{
			name: "linking a product to apikey properly",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID: &apikeyIDs[0], //user0's
				Contracts: []model.AuthorizedContractProducts{
					{
						ContractID: contractIDs[0],
						ProductIDs: []int{productIDs[0]},
					},
				},
			},
			wantStatus: http.StatusCreated,
			wantResp:   "Created",
			wantDBKeys: &dbKeys{
				apiKeyID:          apikeyIDs[0],
				contractProductID: []int{contractProductContentIDs[0]},
			},
		},
		{
			name: "linking contracts in a multiple contract properly, and if ProductIDs field is omitted all products are linked",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID: &apikeyIDs[1], //user1's
				Contracts: []model.AuthorizedContractProducts{
					{
						ContractID: contractIDs[1],
						ProductIDs: []int{ /*productIDs[0]*/ },
					},
					{
						ContractID: contractIDs[2],
					},
				},
			},
			wantStatus: http.StatusCreated,
			wantResp:   "Created",
			wantDBKeys: &dbKeys{
				apiKeyID:          apikeyIDs[1],
				contractProductID: []int{contractProductContentIDs[2], contractProductContentIDs[4], contractProductContentIDs[5]},
			},
		},
		{
			name: "the user of the apikey and one of the contract are different",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID: &apikeyIDs[0], //user0's
				Contracts: []model.AuthorizedContractProducts{
					{
						ContractID: contractIDs[1], //user1's
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: fmt.Sprintf("contract %d does not exist or is not yours, or some products in contract %d are wrong", contractIDs[1], contractIDs[1]),
			},
		},
		{
			name: "the user of the apikey and one of the contract are different",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID: &apikeyIDs[0], //user0's
				Contracts: []model.AuthorizedContractProducts{
					{
						ContractID: contractIDs[0],
						ProductIDs: []int{productIDs[0], productIDs[2]},
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: fmt.Sprintf("following product ids is not found in ids linked to contract %d: [%d]", contractIDs[0], productIDs[2]),
			},
		},
		{
			name: "the apikey does not exist",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID: &notExistID, //user0's
				Contracts: []model.AuthorizedContractProducts{
					{
						ContractID: contractIDs[0],
						ProductIDs: []int{productIDs[0]},
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: fmt.Sprintf("apikey not found, id %d", notExistID),
			},
		},
		{
			name: "api key field is missed",
			req: model.PostAPIKeyProductsReq{
				Contracts: []model.AuthorizedContractProducts{
					{
						ContractID: contractIDs[0],
						ProductIDs: []int{productIDs[0]},
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "apikey_id",
						ConstraintType: "required",
						Message:        "required field, but got empty",
					},
				},
			},
		},
		{
			name: "contracts field is missed",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID: &apikeyIDs[0], //user0's
				Contracts: []model.AuthorizedContractProducts{
					{
						ProductIDs: []int{productIDs[0]},
					},
				},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "contracts[0].contract_id",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            0.0,
					},
				},
			},
		},
		{
			name: "contracts field is empty",
			req: model.PostAPIKeyProductsReq{
				ApiKeyID:  &apikeyIDs[0], //user0's
				Contracts: []model.AuthorizedContractProducts{},
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "contracts",
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

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/keys/products", body)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()
			managementapi.PostAPIKeyProducts(w, r)

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
			if tt.wantDBKeys == nil {
				return
			}

			for i, key := range tt.wantDBKeys.contractProductID {
				rows, err := db.Queryx(`SELECT id
					       				FROM apikey_contract_product_authorized WHERE apikey_id=$1 AND contract_product_id=$2`, tt.wantDBKeys.apiKeyID, key)
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
					t.Errorf("cannot get %dth item, apikey_id %d, contract_product_id %d", i, tt.wantDBKeys.apiKeyID, key)
					return
				}

			}

		})
	}

}
