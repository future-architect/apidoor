package managementapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type productAPIContent struct {
	ApiID       int    `db:"api_id"`
	Description string `db:"description"`
}

func TestPostProduct(t *testing.T) {
	if _, err := db.Exec("TRUNCATE product_api_content"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM product"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Exec("TRUNCATE product_api_content")
		db.Exec("DELETE FROM product")
		db.Exec("DELETE FROM apiinfo")
	}()

	// set up api info db
	apiInfoList := model.APIInfoList{
		List: []model.APIInfo{
			{
				Name:        "info1",
				Source:      "info1 company",
				Description: "description 1",
				Thumbnail:   "http://example.com/img1",
				SwaggerURL:  "http://example.com/test",
			},
			{
				Name:        "info2",
				Source:      "info2 company",
				Description: "description 2",
				Thumbnail:   "http://example.com/img2",
				SwaggerURL:  "http://example.com/test",
			},
		},
	}

	notExistID := int(1e9)

	for i, info := range apiInfoList.List {
		stmt, err := db.PrepareNamed(
			`INSERT INTO apiinfo(name, source, description, thumbnail, swagger_url)
			VALUES (:name, :source, :description, :thumbnail, :swagger_url) RETURNING id`)
		if err != nil {
			t.Fatal(err)
		}
		var id int
		stmt.QueryRowx(info).Scan(&id)
		apiInfoList.List[i].ID = id
	}

	tests := []struct {
		name              string
		contentType       string
		req               model.PostProductReq
		wantHTTPStatus    int
		wantResp          interface{}
		wantProductRecord []model.Product
		wantContentRecord []productAPIContent
	}{
		{
			name:        "post product containing multiple APIs",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "product1",
				DisplayName: "product 1",
				Source:      "company 1",
				Description: "product 1 has two APIs",
				Thumbnail:   "http://example.com/test",
				Contents: []model.APIContent{
					{
						ID:          apiInfoList.List[0].ID,
						Description: "first api",
					},
					{
						ID:          apiInfoList.List[1].ID,
						Description: "second api",
					},
				},
				IsAvailable: false,
			},
			wantHTTPStatus: http.StatusCreated,
			wantResp:       "Created",
			wantProductRecord: []model.Product{
				{
					Name:        "product1",
					DisplayName: "product 1",
					Source:      "company 1",
					Description: "product 1 has two APIs",
					Thumbnail:   "http://example.com/test",
				},
			},
			wantContentRecord: []productAPIContent{
				{
					ApiID:       apiInfoList.List[0].ID,
					Description: "first api",
				},
				{
					ApiID:       apiInfoList.List[1].ID,
					Description: "second api",
				},
			},
		},
		{
			name:        "post product when some optional fields are omitted",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "product2",
				Source:      "company 2",
				Description: "product 2 has one API",
				Thumbnail:   "http://example.com/test",
				Contents: []model.APIContent{
					{
						ID: apiInfoList.List[0].ID,
					},
				},
				IsAvailable: true,
			},
			wantHTTPStatus: http.StatusCreated,
			wantResp:       "Created",
			wantProductRecord: []model.Product{
				{
					Name:            "product2",
					Source:          "company 2",
					Description:     "product 2 has one API",
					Thumbnail:       "http://example.com/test",
					IsAvailableCode: 1,
				},
			},
			wantContentRecord: []productAPIContent{
				{
					ApiID: apiInfoList.List[0].ID,
				},
			},
		},
		{
			name:        "post product containing no API is allowed",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "product3",
				Source:      "company 3",
				Description: "product 3 has no API",
				Thumbnail:   "http://example.com/test",
				Contents:    []model.APIContent{},
				IsAvailable: false,
			},
			wantHTTPStatus: http.StatusCreated,
			wantResp:       "Created",
			wantProductRecord: []model.Product{
				{
					Name:        "product3",
					Source:      "company 3",
					Description: "product 3 has no API",
					Thumbnail:   "http://example.com/test",
				},
			},
			wantContentRecord: []productAPIContent{},
		},
		{
			name:        "api id not existing is contained",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "product91",
				DisplayName: "product 91",
				Source:      "company 91",
				Description: "product 91 is wrong",
				Thumbnail:   "http://example.com/test",
				Contents: []model.APIContent{
					{
						ID:          notExistID,
						Description: "api not existing",
					},
				},
				IsAvailable: false,
			},
			wantHTTPStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message:          fmt.Sprintf("api_id %d does not exist", notExistID),
				ValidationErrors: nil,
			},
			wantProductRecord: []model.Product{},
			wantContentRecord: []productAPIContent{},
		},
		{
			name:        "input validation is failed",
			contentType: "application/json",
			req: model.PostProductReq{
				Name:        "product 92",
				DisplayName: "product 92",
				Description: "product 92 is wrong",
				Thumbnail:   "http://example.com/test",
				Contents:    []model.APIContent{},
				IsAvailable: false,
			},
			wantHTTPStatus: http.StatusBadRequest,
			wantResp: validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "name",
						ConstraintType: "alphanum",
						Message:        "input value, product 92, does not satisfy the format, alphanum",
						Got:            "product 92",
					},
					{
						Field:          "source",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
			wantProductRecord: []model.Product{},
			wantContentRecord: []productAPIContent{},
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

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/products", body)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			managementapi.PostProduct(w, r)

			rw := w.Result()

			resp, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Errorf("read response body error: %v", err)
				return
			}

			if rw.StatusCode != tt.wantHTTPStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantHTTPStatus)
			}

			switch tt.wantResp.(type) {
			case string:
				if tt.wantResp != string(resp) {
					t.Errorf("wrong reponse body: got %s, want %s", resp, tt.wantResp)
				}
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not unsupported")
			}

			// db check
			rows, err := db.Queryx(`SELECT id, name, display_name, source, description, thumbnail, is_available
       				FROM product WHERE name=$1`, tt.req.Name)
			//rows, err := db.Queryx("SELECT * from apiinfo WHERE name=$1", tt.req.Name)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}
			productList := make([]model.Product, 0)
			for rows.Next() {
				var row model.Product
				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}
				productList = append(productList, row)
			}
			if diff := cmp.Diff(tt.wantProductRecord, productList,
				cmpopts.IgnoreFields(model.Product{}, "ID")); diff != "" {
				t.Errorf("db get list of product responce differs:\n %v", diff)
			}

			if len(productList) == 0 {
				return
			}
			productID := productList[0].ID

			rows, err = db.Queryx("SELECT api_id, description from product_api_content WHERE product_id=$1", productID)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}
			contentList := make([]productAPIContent, 0)
			for rows.Next() {
				var row productAPIContent
				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}
				contentList = append(contentList, row)
			}
			if diff := cmp.Diff(tt.wantContentRecord, contentList); diff != "" {
				t.Errorf("db get list of content responce differs:\n %v", diff)
			}
		})
	}
}
