package managementapi

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
)

// SearchProducts godoc
// @Summary search for products
// @Description Get list of APIs and its information
// @produce json
// @Param q query string true "search query words (split words by '.', ex: 'foo.bar'). If containing multiple words, items which match the all search words return"
// @Param target_fields query string false "search target fields. You can choose field(s) from 'all' (represents searching all fields), 'name', 'description', or 'source'. (if there are multiple target fields, split target by '.', ex: 'name.source')" default(all)
// @Param pattern_match query string false "pattern match, chosen from 'exact' or 'partial'" Enums(exact, partial) default(partial)
// @Param limit query int false "the maximum number of results" default(50) minimum(1) maximum(100)
// @Param offset query int false "the starting point for the result set" default(0)
// @Success 200 {object} SearchProductsResp
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /products/search [get]
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("parse param error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	var req SearchProductsReq
	if err := schemaDecoder.Decode(&req, r.Form); err != nil {
		log.Printf("parse query param error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	params, err := req.CreateParams()
	if err != nil {
		log.Printf("validate query param error: %v", err)
		http.Error(w, "param validation error", http.StatusBadRequest)
		return
	}

	respBody, err := db.searchProducts(r.Context(), params)
	if err != nil {
		log.Printf("search products db error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
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
