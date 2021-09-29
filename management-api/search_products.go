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
// @Param q query string true "search query words (split words by '.', ex: 'foo.bar')"
// @Param target_fields query string false "search target fields. You can choose field(s) from 'all' (represents searching all fields), 'name', 'description', or 'source'. (if there are multiple target fields, split target by '.', ex: 'name.source')" default(all)
// @Param pattern_match query string false "pattern match, chosen from 'exact' or 'partial'" Enums(exact, partial) default(partial)
// @Param limit query int false "the maximum number of results" default(50) minimum(1) maximum(100)
// @Param offset query int false "the starting point for the result set" default(0)
// @Success 200 {object} SearchProductsResp
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
