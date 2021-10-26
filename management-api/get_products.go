package managementapi

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// GetProducts godoc
// @Summary Get list of products
// @Description Get list of APIs and its information
// @produce json
// @Success 200 {object} Products
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {
	list, err := db.getProducts(r.Context())
	if err != nil {
		log.Printf("execute get products from db error: %v", err)
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(Products{Products: list})
	if err != nil {
		log.Print("error occurs while reading response")
		http.Error(w, "error occur in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
