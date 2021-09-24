package managementapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// SearchProducts godoc
// @Summary search for products
// @Description Get list of APIs and its information
// @produce json
// @Param name query string false "search words for an API name attribute"
// @Param is_name_exact query boolean false "whether API names are searched by exact match, else searched by partial match" default(false)
// @Param source query string false "search words for an API source attribute"
// @Param is_source_exact query boolean false "whether API names are searched by exact match, else searched by partial match" default(false)
// @Param description query string false "search words for a description attribute by partial match"
// @Param keyword query string false "search words for all attributes"
// @Success 200 {object} Products
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /products/search [get]
func SearchProducts(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("parse param error: %v", err)
		http.Error(w, "parse param error", http.StatusInternalServerError)
		return
	}

	var req SearchProductsReq
	if err := schemaDecoder.Decode(&req, r.Form); err != nil {
		log.Printf("parse query param error: %v", err)
		http.Error(w, "parse query param error", http.StatusInternalServerError)
		return
	}

	rows, err := db.NamedQueryContext(r.Context(), `SELECT * FROM apiinfo WHERE 1=1
AND CASE :name
	WHEN '' THEN 1=1
	ELSE CASE :is_name_exact
		WHEN 'true' THEN name=:name
		ELSE name LIKE concat('%', cast(:name as text), '%')
	END
END
AND CASE :source
	WHEN '' THEN 1=1
	ELSE CASE :is_source_exact
		WHEN 'true' THEN source=:source
		ELSE source LIKE concat('%', cast(:source as text), '%')
	END
END
AND CASE :description
	WHEN '' THEN 1=1
	ELSE description LIKE concat('%', cast(:description as text), '%')
END
AND CASE :keyword
	WHEN '' THEN 1=1
	ELSE name || ' ' || source || ' ' || description LIKE concat('%', cast(:keyword as text), '%')
END
`, req)

	if err != nil {
		log.Printf("error occurs while running query: %v", err)
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}

	list := []Product{}
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
