package managementapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/model"
	swaggerparser "github.com/future-architect/apidoor/managementapi/swagger-parser"
	"github.com/future-architect/apidoor/managementapi/usecase"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPostProduct(t *testing.T) {
	dbType := managementapi.GetAPIDBType(t)
	if dbType != managementapi.DYNAMO {
		log.Println("this test is valid when dynamodb is used, skip")
		return
	}

	managementapi.Setup(t,
		`aws dynamodb --profile local --endpoint-url http://localhost:4566 create-table --cli-input-json file://../dynamo_table/swagger_table.json`,
	)
	t.Cleanup(func() {
		managementapi.Teardown(t,
			`aws dynamodb --profile local --endpoint-url http://localhost:4566 delete-table --table swagger`,
		)
	})

	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}

	usecase.Parser = swaggerparser.NewParser(swaggerparser.TestFetcher{})

	tests := []struct {
		name           string
		contentType    string
		req            model.PostProductReq
		wantHttpStatus int
		//wantRecord は期待されるDB作成レコードの値、idは比較対象外
		wantRecords []model.Product
		wantResp    interface{}
	}{
		{
			name:        "productを登録できる",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "Awesome API",
				Source:      "Company1",
				DisplayName: "display1",
				Description: "provide fantastic product.",
				Thumbnail:   "http://example.com/api.awesome",
				SwaggerURL:  "http://api.example.com/v2/swagger.json",
			},
			wantHttpStatus: http.StatusCreated,
			wantRecords: []model.Product{
				{
					Name:        "Awesome API",
					Source:      "Company1",
					DisplayName: "display1",
					Description: "provide fantastic product.",
					Thumbnail:   "http://example.com/api.awesome",
					BasePath:    "/sample_gateway",
					SwaggerURL:  "http://api.example.com/v2/swagger.json",
				},
			},
			wantResp: model.Product{
				Name:        "Awesome API",
				Source:      "Company1",
				DisplayName: "display1",
				Description: "provide fantastic product.",
				Thumbnail:   "http://example.com/api.awesome",
				BasePath:    "/sample_gateway",
				SwaggerURL:  "http://api.example.com/v2/swagger.json",
			},
		},
		{
			name:        "Fieldに空文字列がある場合は登録できない",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "",
				Source:      "Company2",
				DisplayName: "display2",
				Description: "provide fantastic product.",
				Thumbnail:   "http://example.com/api.awesome",
				SwaggerURL:  "http://api.example.com/v2/swagger.json",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []model.Product{},
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "name",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
		},
		{
			name:        "Content-Typeがapplication/json以外の場合は登録できない",
			contentType: "text/plain",
			req: model.PostProductReq{
				Name:        "wrong content-type",
				Source:      "Company3",
				Description: "provide fantastic product.",
				Thumbnail:   "http://example.com/api.awesome",
				SwaggerURL:  "http://api.example.com/v2/swagger.json",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []model.Product{},
			wantResp: validator.BadRequestResp{
				Message: `unexpected request Content-Type, it must be "application/json"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("create request body error: %v", err)
				return
			}
			body := bytes.NewReader(bodyBytes)

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/api", body)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			managementapi.PostProduct(w, r)

			rw := w.Result()

			resp, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rw.StatusCode != tt.wantHttpStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantHttpStatus)
			}

			rows, err := db.Queryx("SELECT * from product WHERE source=$1", tt.req.Source)
			if err != nil {
				t.Errorf("db get product error: %v", err)
				return
			}

			list := []model.Product{}
			for rows.Next() {
				var row model.Product

				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}

				list = append(list, row)
			}

			if diff := cmp.Diff(tt.wantRecords, list, cmpopts.IgnoreFields(model.Product{}, "ID", "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("db get list of product responce differs:\n %v", diff)
			}
			switch tt.wantResp.(type) {
			case string:
				if tt.wantResp != string(resp) {
					t.Errorf("response body is not %s, got %s", tt.wantResp, string(resp))
				}
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			case model.Product:
				want := tt.wantResp.(model.Product)
				testProduct(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not unsupported")
			}
		})
	}

	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}

}

func testProduct(t *testing.T, want *model.Product, got []byte) {
	t.Helper()
	var gotBody model.Product
	if err := json.Unmarshal(got, &gotBody); err != nil {
		t.Errorf("parsing body as product failed: %v\ngot: %v", err, string(got))
		return
	}

	if diff := cmp.Diff(gotBody, *want, cmpopts.IgnoreFields(model.Product{}, "ID", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Errorf("product differs:\n%v", diff)
	}
}
