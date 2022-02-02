package managementapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/future-architect/apidoor/managementapi"
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
	if _, err := db.Exec("TRUNCATE apiinfo"); err != nil {
		t.Fatal(err)
	}
	defer db.Exec("TRUNCATE apiinfo")

	// set up api info db
	apiInfoList := managementapi.APIInfoList{
		List: []managementapi.APIInfo{
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
		req               managementapi.PostProductReq
		wantHTTPStatus    int
		wantResp          interface{}
		wantProductRecord []managementapi.Product
		wantContentRecord []productAPIContent
	}{
		{
			name:        "post product containing multiple APIs",
			contentType: "application/json",
			req: managementapi.PostProductReq{
				Name:        "product1",
				DisplayName: "product 1",
				Source:      "company 1",
				Description: "product 1 has two APIs",
				Thumbnail:   "http://example.com/test",
				Contents: []managementapi.APIContent{
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
			wantProductRecord: []managementapi.Product{
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
			name:        "api id not existing is contained",
			contentType: "application/json",
			req: managementapi.PostProductReq{
				Name:        "product9",
				DisplayName: "product 9",
				Source:      "company 9",
				Description: "product 9 is wrong",
				Thumbnail:   "http://example.com/test",
				Contents: []managementapi.APIContent{
					{
						ID:          notExistID,
						Description: "api not existing",
					},
				},
				IsAvailable: false,
			},
			wantHTTPStatus: http.StatusBadRequest,
			wantResp: managementapi.BadRequestResp{
				Message:          fmt.Sprintf("api_id %d does not exist", notExistID),
				ValidationErrors: nil,
			},
			wantProductRecord: []managementapi.Product{},
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
			case managementapi.BadRequestResp:
				want := tt.wantResp.(managementapi.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not unsupported")
			}

			// db check
			rows, err := db.Queryx("SELECT * from apiinfo WHERE name=$1", tt.req.Name)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}
			productList := make([]managementapi.Product, 0)
			for rows.Next() {
				var row managementapi.Product
				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}
				productList = append(productList, row)
			}
			if diff := cmp.Diff(tt.wantProductRecord, productList,
				cmpopts.IgnoreFields(managementapi.Product{}, "ID")); diff != "" {
				t.Errorf("db get list of product responce differs:\n %v", diff)
			}

			if len(productList) == 0 {
				return
			}
			productID := productList[0].ID

			rows, err = db.Queryx("SELECT * from product_api_content WHERE product_id=$1", productID)
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
