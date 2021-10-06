package managementapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
		http.Error(w, "parse param error", http.StatusInternalServerError)
		return
	}

	var req SearchProductsReq
	if err := schemaDecoder.Decode(&req, r.Form); err != nil {
		log.Printf("parse query param error: %v", err)
		http.Error(w, "parse query param error", http.StatusInternalServerError)
		return
	}

	params, err := req.CreateParams()
	if err != nil {
		log.Printf("validate query param error: %v", err)
		http.Error(w, "param validation error", http.StatusBadRequest)
		return
	}

	/*
		create SQL statement, such like the following statement:
			SELECT *, COUNT(*) OVER()
			FROM (
				SELECT * FROM apiinfo
				WHERE name LIKE concat('%', cast(:q1 as text), '%')
				UNION
				SELECT * FROM apiinfo
				WHERE source LIKE concat('%', cast(:q1 as text), '%')
			) as T1
			INNER JOIN (
				SELECT * FROM apiinfo
				WHERE name LIKE concat('%', cast(:q2 as text), '%')
				UNION SELECT * FROM apiinfo
				WHERE source LIKE concat('%', cast(:q2 as text), '%')
			) as T2
			on T1.id = T2.id
			ORDER BY T1.id LIMIT :limit OFFSET :offset
	*/
	subQueries := make([]string, len(params.Q))
	for i := range params.Q {
		right := ""
		if params.PatternMatch == "partial" {
			right = fmt.Sprintf("LIKE concat('%%', cast(:q%d as text), '%%')", i+1)
		} else {
			right = fmt.Sprintf("= :q%d", i+1)
		}
		subSubQueries := make([]string, len(params.TargetFields))
		for i, v := range params.TargetFields {
			subSubQueries[i] = fmt.Sprintf("SELECT * FROM apiinfo WHERE %s %s", v, right)
		}
		if i == 0 {
			subQueries[i] = fmt.Sprintf("FROM ( %s ) as T1", strings.Join(subSubQueries, " UNION "))
		} else {
			subQueries[i] = fmt.Sprintf("INNER JOIN ( %s ) as T%d on T1.id = T%d.id",
				strings.Join(subSubQueries, " UNION "), i+1, i+1)
		}
	}
	query := fmt.Sprintf("SELECT *, COUNT(*) OVER() %s ORDER BY T1.id LIMIT :limit OFFSET :offset",
		strings.Join(subQueries, " "))
	targetValues := make(map[string]interface{}, len(params.Q)+2)
	for i, q := range params.Q {
		key := fmt.Sprintf("q%d", i+1)
		targetValues[key] = q
	}
	targetValues["limit"] = params.Limit
	targetValues["offset"] = params.Offset

	rows, err := db.NamedQueryContext(r.Context(), query, targetValues)

	if err != nil {
		log.Printf("error occurs while running query: %v", err)
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}

	list := []Product{}
	count := 0
	for rows.Next() {
		var row SearchProductsResult

		if err := rows.StructScan(&row); err != nil {
			log.Print("error occurs while reading row")
			http.Error(w, "error occurs in database", http.StatusInternalServerError)
			return
		}

		list = append(list, row.Product)
		count = row.Count
	}

	metaData := SearchProductsMetaData{
		ResultSet: ResultSet{
			Count:  count,
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	}

	res, err := json.Marshal(SearchProductsResp{
		Products:               list,
		SearchProductsMetaData: metaData,
	})
	if err != nil {
		log.Print("error occurs while reading response")
		http.Error(w, "error occur in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
