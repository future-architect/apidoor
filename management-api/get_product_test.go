package managementapi_test

import (
	"encoding/json"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/model"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

func init() {
	dbDriver := os.Getenv("DATABASE_DRIVER")
	dbSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSLMODE"))

	var err error
	if db, err = sqlx.Open(dbDriver, dbSource); err != nil {
		log.Fatalf("db connection error: %v", err)
	}
}

func TestGetProducts(t *testing.T) {
	// insert data for test
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}

	var data = []model.Product{
		{
			Name:        "Awesome API",
			Source:      "Nice Company",
			DisplayName: "display1",
			Description: "provide fantastic product.",
			Thumbnail:   "test.com/img/123",
			BasePath:    "test",
			SwaggerURL:  "example.com/api/awesome",
		},
		{
			Name:        "Awesome API v2",
			Source:      "Nice Company",
			DisplayName: "display2",
			Description: "provide special product.",
			Thumbnail:   "test.com/img/456",
			BasePath:    "test",
			SwaggerURL:  "example.com/api/v2/awesome",
		},
	}

	q := `
	INSERT INTO
		product(name, source, display_name, description, thumbnail, swagger_url, base_path, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)
	`
	for _, d := range data {
		if _, err := db.Exec(q, d.Name, d.Source, d.DisplayName, d.Description, d.Thumbnail, d.SwaggerURL, d.BasePath); err != nil {
			t.Fatal(err)
		}
	}

	// test if GetProducts works correctly
	r := httptest.NewRequest(http.MethodGet, "localhost:3000/api", nil)
	w := httptest.NewRecorder()
	managementapi.GetProducts(w, r)

	rw := w.Result()
	defer rw.Body.Close()
	body, err := io.ReadAll(rw.Body)
	if err != nil {
		t.Fatal(err)
	}

	var res model.ProductList
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(data, res.List, cmpopts.IgnoreFields(model.Product{}, "ID", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}

	// reset database
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}
}
