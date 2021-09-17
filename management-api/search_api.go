package managementapi

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

// SearchAPI godoc
// @Summary search for products
// @Description Get list of APIs and its information
// @produce json
// @Param name query string false "search words for an API name attribute"
// @Param is_name_partial query boolean false "whether API names are searched by partial match, else searched by exact match" default(true)
// @Param source query string false "search words for an API source attribute"
// @Param is_source_partial query boolean false "whether API names are searched by partial match, else searched by exact match" default(true)
// @Param description query string false "search words for a description attribute by partial match"
// @Param keyword query string false "search words for all attributes"
// @Success 200 {object} Products
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /products/search [get]
func SearchAPI(w http.ResponseWriter, r *http.Request) {
	descriptionParam := r.URL.Query().Get("description")
	if descriptionParam == "" {
		log.Print("insufficient query parameters")
		http.Error(w, "query parameter, description, missing", http.StatusBadRequest)
	}

	db, err := sqlx.Open(os.Getenv("DATABASE_DRIVER"),
		"host="+os.Getenv("DATABASE_HOST")+" "+
			"port="+os.Getenv("DATABASE_PORT")+" "+
			"user="+os.Getenv("DATABASE_USER")+" "+
			"password="+os.Getenv("DATABASE_PASSWORD")+" "+
			"dbname="+os.Getenv("DATABASE_NAME")+" "+
			"sslmode="+os.Getenv("DATABASE_SSLMODE"))
	if err != nil {
		log.Print("error occurs in database")
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Queryx("SELECT * FROM apiinfo WHERE description=$1", descriptionParam)
	if err != nil {
		log.Print("error occurs while running query")
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
