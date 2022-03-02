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

func TestSearchAPIInfo(t *testing.T) {
	// insert data for test
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

	var data = []model.APIInfo{
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
		params     model.SearchAPIInfoReq
		wantStatus int
		wantResp   interface{} // *managementapi.SearchAPIInfoResp
	}{
		{
			name: "完全一致の検索ができる",
			params: model.SearchAPIInfoReq{
				Q:            "Awesome API",
				PatternMatch: "exact",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchAPIInfoResp{
				APIList: []model.APIInfo{
					{
						Name:        "Awesome API",
						Source:      "Nice Company",
						Description: "provide fantastic information.",
						Thumbnail:   "test.com/img/aaa",
						SwaggerURL:  "example.com/api/awesome",
					},
				},
				SearchAPIInfoMetaData: model.SearchAPIInfoMetaData{
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
			params: model.SearchAPIInfoReq{
				Q: "Awesome API",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchAPIInfoResp{
				APIList: []model.APIInfo{
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
				SearchAPIInfoMetaData: model.SearchAPIInfoMetaData{
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
			params: model.SearchAPIInfoReq{
				Q: "Search.example%2ecom",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchAPIInfoResp{
				APIList: []model.APIInfo{
					{
						Name:        "Search API",
						Source:      "Great Company",
						Description: "search for example.com.",
						Thumbnail:   "test.com/img/ddd",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchAPIInfoMetaData: model.SearchAPIInfoMetaData{
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
			params: model.SearchAPIInfoReq{
				Q:            "Great",
				TargetFields: "source.description",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchAPIInfoResp{
				APIList: []model.APIInfo{
					{
						Name:        "Search API",
						Source:      "Great Company",
						Description: "search for example.com.",
						Thumbnail:   "test.com/img/ddd",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchAPIInfoMetaData: model.SearchAPIInfoMetaData{
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
			params: model.SearchAPIInfoReq{
				Q:            "special",
				PatternMatch: "partial",
				Offset:       1,
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchAPIInfoResp{
				APIList: []model.APIInfo{
					{
						Name:        "Great API",
						Source:      "Nice Company",
						Description: "provide special information.",
						Thumbnail:   "test.com/img/ccc",
						SwaggerURL:  "example.com/api/great",
					},
				},
				SearchAPIInfoMetaData: model.SearchAPIInfoMetaData{
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
			params: model.SearchAPIInfoReq{
				Q: "not exist",
			},
			wantStatus: http.StatusOK,
			wantResp: model.SearchAPIInfoResp{
				APIList: []model.APIInfo{},
				SearchAPIInfoMetaData: model.SearchAPIInfoMetaData{
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
			params: model.SearchAPIInfoReq{
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
			params: model.SearchAPIInfoReq{
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
			params: model.SearchAPIInfoReq{
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
			r := httptest.NewRequest(http.MethodGet, "localhost:3000/api/search?"+form.Encode(), nil)
			w := httptest.NewRecorder()
			managementapi.SearchAPIInfo(w, r)

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
			case model.SearchAPIInfoResp:
				testCompareBodyAsSearchAPIInfoResp(t, tt.wantResp.(model.SearchAPIInfoResp), body)
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, body)
			case string:
				testCompareBodyAsString(t, tt.wantResp.(string), body)
			default:
				t.Error("wantResp is neither SearchAPIInfoResp, BadRequestResp nor string")
			}

		})
	}
}

func testCompareBodyAsSearchAPIInfoResp(t *testing.T, want model.SearchAPIInfoResp, got []byte) {
	var respBody model.SearchAPIInfoResp
	if err := json.Unmarshal(got, &respBody); err != nil {
		t.Errorf("parse body as search api info response error: %v", err)
		return
	}
	if diff := cmp.Diff(want, respBody, cmpopts.IgnoreFields(model.APIInfo{}, "ID")); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}
}

func testCompareBodyAsString(t *testing.T, want string, got []byte) {
	if diff := cmp.Diff(want, string(got)); diff != "" {
		t.Errorf("unexpected response: differs=\n%v", diff)
	}
}
