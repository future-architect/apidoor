package managementapi_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"managementapi"

	_ "github.com/lib/pq"
)

func TestGetProducts(t *testing.T) {
	db := managementapi.DB
	// insert data for test

	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

	var data = []managementapi.Product{
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

	if diff := cmp.Diff(data, res.Products, cmpopts.IgnoreFields(managementapi.Product{}, "ID")); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}

	// reset database
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}
}
