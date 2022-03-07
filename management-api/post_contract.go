package managementapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"log"
	"net/http"
)

// PostContract godoc
// @Summary Post a product
// @Description Post an API product
// @produce json
// @Param product body model.PostContractReq true "contract definition"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /contract [post]
func PostContract(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		resp := validator.NewBadRequestResp(`unexpected request Content-Type, it must be "application/json"`)
		if respBytes, err := json.Marshal(resp); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(respBytes)
		} else {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	body := new(bytes.Buffer)
	io.Copy(body, r.Body)

	var req model.PostContractReq
	if err := json.Unmarshal(body.Bytes(), &req); err != nil {
		if errors.Is(err, model.UnmarshalJsonErr) {
			log.Printf("failed to parse json body: %v", err)
			resp := validator.NewBadRequestResp(model.UnmarshalJsonErr.Error())
			if respBytes, err := json.Marshal(resp); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if ve, ok := err.(validator.ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if respBytes, err := json.Marshal(ve.ToBadRequestResp()); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			// unreachable
			log.Printf("unexpected error returned: %v", err)
			http.Error(w, fmt.Sprintf("server error"), http.StatusInternalServerError)
		}
		return
	}

	userID, err := fetchUserID(r.Context(), req.UserAccountId)
	if err != nil {
		log.Printf("fetch user id error: %v", err)
		if errors.Is(err, ErrNotFound) {
			br := validator.BadRequestResp{
				Message: fmt.Sprintf("account_id %s does not exist", req.UserAccountId),
			}
			if respBytes, err := json.Marshal(br); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	productID, err := fetchProductID(r.Context(), req.ProductName)
	if err != nil {
		log.Printf("fetch product id error: %v", err)
		if errors.Is(err, ErrNotFound) {
			br := validator.BadRequestResp{
				Message: fmt.Sprintf("product_name %s does not exist", req.ProductName),
			}
			if respBytes, err := json.Marshal(br); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	contract := model.Contract{
		UserID:    userID,
		ProductID: productID,
	}

	if err := db.postContract(r.Context(), &contract); err != nil {
		log.Printf("db insert contract error: %v", err)
		// this error occurs the api_user or the product is deleted after fetching its id
		if constraintErr, ok := err.(*dbConstraintErr); ok {
			var errMsg string
			switch constraintErr.field {
			case "user_id":
				errMsg = fmt.Sprintf("account_id %s does not exist", req.UserAccountId)
			case "contract_id":
				errMsg = fmt.Sprintf("product_name %s does not exist", req.ProductName)
			}
			br := validator.BadRequestResp{
				Message: errMsg,
			}
			if respBytes, err := json.Marshal(br); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(respBytes)
			} else {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}

func fetchUserID(ctx context.Context, accountID string) (int, error) {
	// TODO　契約する権限の確認 (管理者ユーザと被管理ユーザのテーブルを作って用いる?)
	user, err := db.fetchUser(ctx, accountID)
	if err != nil {
		return 0, fmt.Errorf("fetch user id db error: %w", err)
	}

	return user.ID, nil
}

func fetchProductID(ctx context.Context, productName string) (int, error) {
	product, err := db.fetchProduct(ctx, productName)
	if err != nil {
		return 0, fmt.Errorf("fetch product id db error: %w", err)
	}

	return product.ID, nil
}
