package managementapi

import (
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/usecase"
	"github.com/future-architect/apidoor/managementapi/validator"
	"log"
	"net/http"
)

// SearchProduct godoc
// @Summary search for products
// @Description search products
// @produce json
// @Param q query string true "search query words (split words by '.', ex: 'foo.bar'). If containing multiple words, items which match the all search words return"
// @Param target_fields query string false "search target fields. You can choose field(s) from 'all' (represents searching all fields), 'name', 'description', or 'source'. (if there are multiple target fields, split target by '.', ex: 'name.source')" default(all)
// @Param pattern_match query string false "pattern match, chosen from 'exact' or 'partial'" Enums(exact, partial) default(partial)
// @Param limit query int false "the maximum number of results" default(50) minimum(1) maximum(100)
// @Param offset query int false "the starting point for the result set" default(0)
// @Success 200 {object} model.SearchProductResp
// @Failure 400 {object} validator.BadRequestResp
// @Failure 500 {string} string
// @Router /products/search [get]
func SearchProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("parse param error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var req model.SearchProductReq
	if err := model.SchemaDecoder.Decode(&req, r.Form); err != nil {
		log.Printf("parse query param error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	params, err := req.CreateParams()
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			log.Printf("input validation failed:\n%v", err)
			writeErrResponse(w, ve)
		} else {
			log.Printf("validate query param error: %v", err)
			writeErrResponse(w, usecase.NewClientError(errors.New("param validation error")))
		}
		return
	}

	respBody, err := usecase.SearchProduct(r.Context(), params)
	if err != nil {
		writeErrResponse(w, err)
		return
	}

	res, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("create json response error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
