package managementapi_test

import (
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/schema"
)

func TestSearchProduct(t *testing.T) {
	// insert data for test
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}

	var data = []model.Product{
		{
			Name:        "Awesome API",
			Source:      "Nice Company",
			DisplayName: "display",
			Description: "provide fantastic product.",
			Thumbnail:   "test.com/img/aaa",
			SwaggerURL:  "example.com/api/awesome",
		},
		{
			Name:        "Awesome API v2",
			Source:      "Very Nice Company",
			DisplayName: "display",
			Description: "provide special product.",
			Thumbnail:   "test.com/img/bbb",
			SwaggerURL:  "example.com/api/v2/awesome",
		},
		{
			Name:        "Great API",
			Source:      "Nice Company",
			DisplayName: "display",
			Description: "provide special product.",
			Thumbnail:   "test.com/img/ccc",
			SwaggerURL:  "example.com/api/great",
		},
		{
			Name:        "Search API",
			Source:      "Great Company",
			DisplayName: "display",
			Description: "search for example.com.",
			Thumbnail:   "test.com/img/ddd",
			SwaggerURL:  "example.com/api/great",
		},
		{
			Name:        "Search API2",
			Source:      "Good Company",
			DisplayName: "display",
			Description: "search for example.net.",
			Thumbnail:   "test.com/img/ddd",
			SwaggerURL:  "example.com/api/great",
		},
	}

	q := `
	INSERT INTO
		product(name, source, display_name, description, thumbnail, swagger_url, base_path, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, current_timestamp, current_timestamp)
	`
	for _, d := range data {
		if _, err := db.Exec(q, d.Name, d.Source, d.DisplayName, d.Description, d.Thumbnail, d.SwaggerURL, "/foo"); err != nil {
			t.Fatal(err)
		}
	}

	encoder := schema.NewEncoder()

	tests := []struct {
		name       string
		params     model.SearchProductReq
		wantStatus int
		wantResp   interface{} // *management-api.SearchProductResp
	}{
		{
			name: "完全一致の検索ができる",
			params: model.SearchProductReq{
				Q:            "Awesome API",
				PatternMatch: "exact",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchProductResp{
				ProductList: []model.Product{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic product.",
						Thumbnail:   "test.com/img/aaa",
						SwaggerURL:  "example.com/api/awesome",
					},
				},
				SearchProductMetaData: model.SearchProductMetaData{
					ResultSet: model.ResultSet{
						Count:  1,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "部分一致の検索ができる(pattern matchは省略可能)",
			params: model.SearchProductReq{
				Q: "Awesome API",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchProductResp{
				ProductList: []model.Product{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic product.",
						Thumbnail:   "test.com/img/aaa",
						SwaggerURL:  "example.com/api/awesome",
					},
					{
						Name:        "Awesome API v2",
						Source:      "Very Nice Company",
						Description: "provide special product.",
						Thumbnail:   "test.com/img/bbb",
						SwaggerURL:  "example.com/api/v2/awesome",
					},
				},
				SearchProductMetaData: model.SearchProductMetaData{
					ResultSet: model.ResultSet{
						Count:  2,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "複数キーワードであり、また、パーセントエンコーディングを持つキーワードを含む部分一致検索ができる",
			params: model.SearchProductReq{
				Q: "Search.example%2ecom",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchProductResp{
				ProductList: []model.Product{
					{
						Name:        "Search API",
						Source:      "Great Company",
						Description: "search for example.com.",
						Thumbnail:   "test.com/img/ddd",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchProductMetaData: model.SearchProductMetaData{
					ResultSet: model.ResultSet{
						Count:  1,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "フィールドを指定して検索ができる",
			params: model.SearchProductReq{
				Q:            "Great",
				TargetFields: "source.description",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchProductResp{
				ProductList: []model.Product{
					{
						Name:        "Search API",
						Source:      "Great Company",
						Description: "search for example.com.",
						Thumbnail:   "test.com/img/ddd",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchProductMetaData: model.SearchProductMetaData{
					ResultSet: model.ResultSet{
						Count:  1,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "limitで件数を制限し、offsetで開始位置を指定できる",
			params: model.SearchProductReq{
				Q:            "special",
				PatternMatch: "partial",
				Offset:       1,
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchProductResp{
				ProductList: []model.Product{
					{
						Name:        "Great API",
						Source:      "Nice Company",
						Description: "provide special product.",
						Thumbnail:   "test.com/img/ccc",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchProductMetaData: model.SearchProductMetaData{
					ResultSet: model.ResultSet{
						Count:  2,
						Limit:  50,
						Offset: 1,
					},
				},
			},
		},
		{
			name: "検索結果が0件",
			params: model.SearchProductReq{
				Q: "not exist",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchProductResp{
				ProductList: []model.Product{},
				SearchProductMetaData: model.SearchProductMetaData{
					ResultSet: model.ResultSet{
						Count:  0,
						Limit:  50,
						Offset: 0,
					},
				},
			},
		},
		{
			name: "リクエストパラメータが不正",
			params: model.SearchProductReq{
				Q:            "img",
				TargetFields: "name.thumbnail",
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "target_fields[1]",
						ConstraintType: "enum",
						Message:        "input value is thumbnail, but it must be one of the following values: [all name description source]",
						Enum:           []string{"all", "name", "description", "source"},
						Got:            "thumbnail",
					},
				},
			},
		},
		{
			name: "Qパラメータが未指定、または空文字列",
			params: model.SearchProductReq{
				TargetFields: "name",
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "q",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
		},
		{
			name: "Qパラメータに空文字列が含まれている",
			params: model.SearchProductReq{
				Q:            "Awesome..API",
				TargetFields: "name",
			},
			wantStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "q[1]",
						ConstraintType: "ne",
						Message:        "input value is , but it must be not equal to ",
						Got:            "",
					},
				},
			},
		},
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM product"); err != nil {
			t.Fatal(err)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			if err := encoder.Encode(tt.params, form); err != nil {
				t.Fatalf("encode params error: %v", err)
			}
			r := httptest.NewRequest(http.MethodGet, "localhost:3000/api/search?"+form.Encode(), nil)
			w := httptest.NewRecorder()
			managementapi.SearchProduct(w, r)

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
			case model.SearchProductResp:
				testCompareBodyAsSearchProductResp(t, tt.wantResp.(model.SearchProductResp), body)
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, body)
			case string:
				testCompareBodyAsString(t, tt.wantResp.(string), body)
			default:
				t.Error("wantResp is neither SearchProductResp, BadRequestResp nor string")
			}

		})
	}
}

func testCompareBodyAsSearchProductResp(t *testing.T, want model.SearchProductResp, got []byte) {
	var respBody model.SearchProductResp
	if err := json.Unmarshal(got, &respBody); err != nil {
		t.Errorf("parse body as search product response error: %v", err)
		return
	}
	if diff := cmp.Diff(want, respBody,
		cmpopts.IgnoreFields(model.Product{}, "ID", "DisplayName", "BasePath", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}
}

func testCompareBodyAsString(t *testing.T, want string, got []byte) {
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}
}
