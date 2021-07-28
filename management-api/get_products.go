package managementapi

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// GetProducts godoc
// @Summary Get list of products
// @Description Get list of APIs and its information
// @produce json
// @Success 200 {object} Products
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Queryx("SELECT * from apiinfo")
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
