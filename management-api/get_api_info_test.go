package managementapi_test

import (
	"encoding/json"
	"fmt"
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

func TestGetAPIInfo(t *testing.T) {
	// insert data for test
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

	var data = []managementapi.APIInfo{
		{
			Name:        "Awesome API",
			Source:      "Nice Company",
			Description: "provide fantastic information.",
			Thumbnail:   "test.com/img/123",
			SwaggerURL:  "example.com/api/awesome",
		},
		{
			Name:        "Awesome API v2",
			Source:      "Nice Company",
			Description: "provide special information.",
			Thumbnail:   "test.com/img/456",
			SwaggerURL:  "example.com/api/v2/awesome",
		},
	}

	q := `
	INSERT INTO
		apiinfo(name, source, description, thumbnail, swagger_url)
	VALUES
		($1, $2, $3, $4, $5)
	`
	for _, d := range data {
		if _, err := db.Exec(q, d.Name, d.Source, d.Description, d.Thumbnail, d.SwaggerURL); err != nil {
			t.Fatal(err)
		}
	}

	// test if GetAPIInfo works correctly
	r := httptest.NewRequest(http.MethodGet, "localhost:3000/api", nil)
	w := httptest.NewRecorder()
	managementapi.GetAPIInfo(w, r)

	rw := w.Result()
	defer rw.Body.Close()
	body, err := io.ReadAll(rw.Body)
	if err != nil {
		t.Fatal(err)
	}

	var res managementapi.APIInfoList
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(data, res.List, cmpopts.IgnoreFields(managementapi.APIInfo{}, "ID")); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}

	// reset database
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}
}
