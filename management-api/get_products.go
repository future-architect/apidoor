package managementapi

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// get list of information of APIs from database
func GetProducts(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(os.Getenv("DATABASE_DRIVER"),
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

	rows, err := db.Query("SELECT * from apiinfo")
	if err != nil {
		log.Print("error occurs while running query")
		http.Error(w, "error occurs in database", http.StatusInternalServerError)
		return
	}

	var list []Api
	for rows.Next() {
		var row Api

		if err := rows.Scan(&row.ID, &row.Name, &row.Source, &row.Description, &row.Thumbnail); err != nil {
			log.Print("error occurs while reading row")
			http.Error(w, "error occurs in database", http.StatusInternalServerError)
			return
		}

		list = append(list, row)
	}

	res, err := json.Marshal(list)
	if err != nil {
		log.Print("error occurs while reading response")
		http.Error(w, "error occur in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
