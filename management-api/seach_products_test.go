package managementapi_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/schema"
	"io"
	"managementapi"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSearchProducts(t *testing.T) {
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
			Source:      "Very Nice Company",
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

	encoder := schema.NewEncoder()

	tests := []struct {
		name       string
		params     managementapi.SearchProductsReq
		wantStatus int
		wantResp   managementapi.Products
	}{
		{
			name: "nameで完全一致の検索ができる",
			params: managementapi.SearchProductsReq{
				Name:        "Awesome API",
				IsNameExact: true,
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.Products{
				Products: []managementapi.Product{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic information.",
						Thumbnail:   "test.com/img/123",
						SwaggerURL:  "example.com/api/awesome",
					},
				},
			},
		},
		{
			name: "sourceで部分一致の検索ができる",
			params: managementapi.SearchProductsReq{
				Source:        "Nice",
				IsSourceExact: false,
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.Products{
				Products: []managementapi.Product{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic information.",
						Thumbnail:   "test.com/img/123",
						SwaggerURL:  "example.com/api/awesome",
					},
					{
						Name:        "Awesome API v2",
						Source:      "Very Nice Company",
						Description: "provide special information.",
						Thumbnail:   "test.com/img/456",
						SwaggerURL:  "example.com/api/v2/awesome",
					},
				},
			},
		},
		{
			name: "keywordで部分一致の検索ができる",
			params: managementapi.SearchProductsReq{
				Keyword: "special",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.Products{
				Products: []managementapi.Product{
					{
						Name:        "Awesome API v2",
						Source:      "Very Nice Company",
						Description: "provide special information.",
						Thumbnail:   "test.com/img/456",
						SwaggerURL:  "example.com/api/v2/awesome",
					},
				},
			},
		},
		{
			name: "検索結果が0件",
			params: managementapi.SearchProductsReq{
				Keyword: "not exist",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.Products{
				Products: []managementapi.Product{},
			},
		},
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
			t.Fatal(err)
		}
	}()

	// test if GetProducts works correctly
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			if err := encoder.Encode(tt.params, form); err != nil {
				t.Fatalf("encode params error: %v", err)
			}
			r := httptest.NewRequest(http.MethodGet, "localhost:3000/prouct/search?"+form.Encode(), nil)
			w := httptest.NewRecorder()
			managementapi.SearchProducts(w, r)

			rw := w.Result()
			defer rw.Body.Close()

			if rw.StatusCode != tt.wantStatus {
				t.Errorf("wrong status code: got %d, want %d", rw.StatusCode, tt.wantStatus)
			}

			body, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatal(err)
			}

			var res managementapi.Products
			if err := json.Unmarshal(body, &res); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.wantResp.Products, res.Products, cmpopts.IgnoreFields(managementapi.Product{}, "ID")); diff != "" {
				t.Errorf("unexpected response: differs=\n%v", diff)
			}

			// reset database

		})
	}
}
