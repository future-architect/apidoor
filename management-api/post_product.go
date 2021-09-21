package managementapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// PostProduct godoc
// @Summary Get list of products
// @Description Get list of APIs and its information
// @produce json
// @Param product body PostProductReq true "api information"
// @Success 201 {string} string
// @Failure 400 {string} error
// @Failure 500 {string} error
// @Router /product [post]
func PostProduct(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		log.Print("unexpected request content")
		http.Error(w, "unexpected request content", http.StatusBadRequest)
		return
	}

	var req PostProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to parse json body: %v", err)
		http.Error(w, "failed to parse json body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		log.Printf("invalid body: %v", err)
		http.Error(w, fmt.Sprintf("invalid body"), http.StatusBadRequest)
		return
	}

	_, err := db.NamedExecContext(r.Context(),
		"INSERT INTO apiinfo(name, source, description, thumbnail, swagger_url) VALUES(:name, :source, :description, :thumbnail, :swagger_url)",
		req)
	if err != nil {
		log.Printf("db insert product error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Created")
}
