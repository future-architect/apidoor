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
			Thumbnail:   "test.com/img/aaa",
			SwaggerURL:  "example.com/api/awesome",
		},
		{
			Name:        "Awesome API v2",
			Source:      "Very Nice Company",
			Description: "provide special information.",
			Thumbnail:   "test.com/img/bbb",
			SwaggerURL:  "example.com/api/v2/awesome",
		},
		{
			Name:        "Great API",
			Source:      "Nice Company",
			Description: "provide special information.",
			Thumbnail:   "test.com/img/ccc",
			SwaggerURL:  "example.com/api/great",
		},
		{
			Name:        "Search API",
			Source:      "Great Company",
			Description: "search for example.com.",
			Thumbnail:   "test.com/img/ddd",
			SwaggerURL:  "example.com/api/great",
		},
		{
			Name:        "Search API2",
			Source:      "Good Company",
			Description: "search for example.net.",
			Thumbnail:   "test.com/img/ddd",
			SwaggerURL:  "example.com/api/great",
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
		wantResp   interface{} // *managementapi.SearchProductsResp
	}{
		{
			name: "完全一致の検索ができる",
			params: managementapi.SearchProductsReq{
				Q:            "Awesome API",
				PatternMatch: "exact",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.SearchProductsResp{
				Products: []managementapi.Product{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic information.",
						Thumbnail:   "test.com/img/aaa",
						SwaggerURL:  "example.com/api/awesome",
					},
				},
				SearchProductsMetaData: managementapi.SearchProductsMetaData{
					ResultSet: managementapi.ResultSet{
						Count:  1,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "部分一致の検索ができる(pattern matchは省略可能)",
			params: managementapi.SearchProductsReq{
				Q: "Awesome API",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.SearchProductsResp{
				Products: []managementapi.Product{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic information.",
						Thumbnail:   "test.com/img/aaa",
						SwaggerURL:  "example.com/api/awesome",
					},
					{
						Name:        "Awesome API v2",
						Source:      "Very Nice Company",
						Description: "provide special information.",
						Thumbnail:   "test.com/img/bbb",
						SwaggerURL:  "example.com/api/v2/awesome",
					},
				},
				SearchProductsMetaData: managementapi.SearchProductsMetaData{
					ResultSet: managementapi.ResultSet{
						Count:  2,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "複数キーワードであり、また、パーセントエンコーディングを持つキーワードを含む部分一致検索ができる",
			params: managementapi.SearchProductsReq{
				Q: "Search.example%2ecom",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.SearchProductsResp{
				Products: []managementapi.Product{
					{
						Name:        "Search API",
						Source:      "Great Company",
						Description: "search for example.com.",
						Thumbnail:   "test.com/img/ddd",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchProductsMetaData: managementapi.SearchProductsMetaData{
					ResultSet: managementapi.ResultSet{
						Count:  1,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "フィールドを指定して検索ができる",
			params: managementapi.SearchProductsReq{
				Q:            "Great",
				TargetFields: "source.description",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.SearchProductsResp{
				Products: []managementapi.Product{
					{
						Name:        "Search API",
						Source:      "Great Company",
						Description: "search for example.com.",
						Thumbnail:   "test.com/img/ddd",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchProductsMetaData: managementapi.SearchProductsMetaData{
					ResultSet: managementapi.ResultSet{
						Count:  1,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "limitで件数を制限し、offsetで開始位置を指定できる",
			params: managementapi.SearchProductsReq{
				Q:            "special",
				PatternMatch: "partial",
				Offset:       1,
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.SearchProductsResp{
				Products: []managementapi.Product{
					{
						Name:        "Great API",
						Source:      "Nice Company",
						Description: "provide special information.",
						Thumbnail:   "test.com/img/ccc",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchProductsMetaData: managementapi.SearchProductsMetaData{
					ResultSet: managementapi.ResultSet{
						Count:  2,
						Limit:  50,
						Offset: 1,
					},
				},
			},
		},
		{
			name: "検索結果が0件",
			params: managementapi.SearchProductsReq{
				Q: "not exist",
			},
			wantStatus: http.StatusOK,
			wantResp: managementapi.SearchProductsResp{
				Products: []managementapi.Product{},
				SearchProductsMetaData: managementapi.SearchProductsMetaData{
					ResultSet: managementapi.ResultSet{
						Count:  0,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "リクエストパラメータが不正",
			params: managementapi.SearchProductsReq{
				Q:            "img",
				TargetFields: "name.thumbnail",
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   "param validation error\n",
		},
		{
			name: "Qパラメータが未指定、または空文字列",
			params: managementapi.SearchProductsReq{
				TargetFields: "name",
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   "param validation error\n",
		},
		{
			name: "Qパラメータに空文字列が含まれている",
			params: managementapi.SearchProductsReq{
				Q:            "Awesome..API",
				TargetFields: "name",
			},
			wantStatus: http.StatusBadRequest,
			wantResp:   "param validation error\n",
		},
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
			t.Fatal(err)
		}
	}()

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
			switch tt.wantResp.(type) {
			case managementapi.SearchProductsResp:
				testCompareBodyAsSearchProductResp(t, tt.wantResp.(managementapi.SearchProductsResp), body)
			case string:
				testCompareBodyAsString(t, tt.wantResp.(string), body)
			default:
				t.Error("wantResp is neither SearchProductResp nor string")
			}

		})
	}
}

func testCompareBodyAsSearchProductResp(t *testing.T, want managementapi.SearchProductsResp, got []byte) {
	var respBody managementapi.SearchProductsResp
	if err := json.Unmarshal(got, &respBody); err != nil {
		t.Errorf("parse body as search products response error: %v", err)
		return
	}
	if diff := cmp.Diff(want, respBody, cmpopts.IgnoreFields(managementapi.Product{}, "ID")); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}
}

func testCompareBodyAsString(t *testing.T, want string, got []byte) {
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}
}
