package managementapi

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var (
	//go:embed sql/search_api.sql
	searchAPISQLTemplateStr string
	searchAPISQLTemplate    *template.Template
)

func init() {
	var err error
	searchAPISQLTemplate, err = template.New("search API  SQL template").Parse(searchAPISQLTemplateStr)
	if err != nil {
		log.Fatalf("create searchAPISQL template %v", err)
	}
}

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

	var query bytes.Buffer
	if err := searchAPISQLTemplate.Execute(&query, params); err != nil {
		log.Printf("generate SQL error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	targetValues := make(map[string]interface{}, len(params.Q)+2)
	for i, q := range params.Q {
		key := fmt.Sprintf("q%d", i)
		targetValues[key] = q
	}
	targetValues["limit"] = params.Limit
	targetValues["offset"] = params.Offset

	rows, err := db.NamedQueryContext(r.Context(), query.String(), targetValues)

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
