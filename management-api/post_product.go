package managementapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"log"
	"net/http"
)

// PostProduct godoc
// @Summary Post a product
// @Description Post an API product
// @produce json
// @Param product body model.PostProductReq true "product definition"
// @Success 201 {string} string
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} error
// @Router /products [post]
func PostProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Printf("unexpected request content: %s", r.Header.Get("Content-Type"))
		resp := validator.NewBadRequestResp(`unexpected request Content-Type, it must be "application/json"`)
		if err := resp.WriteResp(w); err != nil {
			log.Printf("write bad request response failed: %v", err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	body := new(bytes.Buffer)
	io.Copy(body, r.Body)

	var req model.PostProductReq
	if err := json.Unmarshal(body.Bytes(), &req); err != nil {
		if errors.Is(err, model.UnmarshalJsonErr) {
			log.Printf("failed to parse json body: %v", err)
			resp := validator.NewBadRequestResp(model.UnmarshalJsonErr.Error())
			if err := resp.WriteResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if ve, ok := err.(validator.ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			if err = ve.ToBadRequestResp().WriteResp(w); err != nil {
				log.Printf("write bad request response failed: %v", err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else {
			// unreachable
			log.Printf("invalid body: %v", err)
			http.Error(w, fmt.Sprintf("invalid body"), http.StatusBadRequest)
		}
		return
	}
	req = req.Convert()

	if err := db.postProduct(r.Context(), &req); err != nil {
		log.Printf("insert product to db failed: %v", err)
		if constraintErr, ok := err.(*dbConstraintErr); ok {
			br := validator.BadRequestResp{
				Message: fmt.Sprintf("api_id %d does not exist", constraintErr.value),
			}
			if err = br.WriteResp(w); err != nil {
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
