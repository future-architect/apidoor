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
	rows, err := DB.Queryx("SELECT * from apiinfo")
	if err != nil {
		log.Printf("error occurs while running query %v", err)
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}

	var list []Product
	for rows.Next() {
		var row Product

		if err := rows.StructScan(&row); err != nil {
			log.Print("error occurs while reading row")
			http.Error(w, "error occurs in database", http.StatusInternalServerError)
			return
		}

		list = append(list, row)
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
