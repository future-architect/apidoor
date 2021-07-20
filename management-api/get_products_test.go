package managementapi_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"managementapi"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestGetProducts(t *testing.T) {
	// insert data for test
	db, err := sqlx.Open(os.Getenv("DATABASE_DRIVER"),
		"host="+os.Getenv("DATABASE_HOST")+" "+
			"port="+os.Getenv("DATABASE_PORT")+" "+
			"user="+os.Getenv("DATABASE_USER")+" "+
			"password="+os.Getenv("DATABASE_PASSWORD")+" "+
			"dbname="+os.Getenv("DATABASE_NAME")+" "+
			"sslmode="+os.Getenv("DATABASE_SSLMODE"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

	var data = []managementapi.Product{
		{
			ID:          3,
			Name:        "Awesome API",
			Source:      "Nice Company",
			Description: "provide fantastic information.",
			Thumbnail:   "test.com/img/123",
		},
		{
			ID:          4,
			Name:        "Awesome API v2",
			Source:      "Nice Company",
			Description: "provide special information.",
			Thumbnail:   "test.com/img/456",
		},
	}

	q := `
	INSERT INTO 
		apiinfo(id, name, source, description, thumbnail)
	VALUES
		($1, $2, $3, $4, $5)
	`
	for _, d := range data {
		if _, err := db.Exec(q, d.ID, d.Name, d.Source, d.Description, d.Thumbnail); err != nil {
			t.Fatal(err)
		}
	}

	// test if GetProducts works correctly
	r := httptest.NewRequest(http.MethodGet, "localhost:3000/products", nil)
	w := httptest.NewRecorder()
	managementapi.GetProducts(w, r)

	rw := w.Result()
	defer rw.Body.Close()
	body, err := io.ReadAll(rw.Body)
	if err != nil {
		t.Fatal(err)
	}

	var res managementapi.Products
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(data, res.Products) {
		t.Fatalf("unexpected response: expected %v, get %v", data, res)
	}

	// reset database
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}
}
