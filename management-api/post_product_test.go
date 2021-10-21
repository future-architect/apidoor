package managementapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPostProduct(t *testing.T) {
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		contentType    string
		req            managementapi.PostProductReq
		wantHttpStatus int
		//wantRecord は期待されるDB作成レコードの値、idは比較対象外
		wantRecords []managementapi.Product
	}{
		{
			name:        "productを登録できる",
			contentType: "application/json",
			req: managementapi.PostProductReq{
				Name:        "Awesome API",
				Source:      "Company1",
				Description: "provide fantastic information.",
				Thumbnail:   "test.com/api.awesome",
				SwaggerURL:  "example.com/api/awesome",
			},
			wantHttpStatus: http.StatusCreated,
			wantRecords: []managementapi.Product{
				{
					Name:        "Awesome API",
					Source:      "Company1",
					Description: "provide fantastic information.",
					Thumbnail:   "test.com/api.awesome",
					SwaggerURL:  "example.com/api/awesome",
				},
			},
		},
		{
			name:        "Fieldに空文字列がある場合は登録できない",
			contentType: "application/json",
			req: managementapi.PostProductReq{
				Name:        "",
				Source:      "Company2",
				Description: "provide fantastic information.",
				Thumbnail:   "test.com/api.awesome",
				SwaggerURL:  "example.com/api/awesome",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []managementapi.Product{},
		},
		{
			name:        "Content-Typeがapplication/json以外の場合は登録できない",
			contentType: "text/plain",
			req: managementapi.PostProductReq{
				Name:        "wrong content-type",
				Source:      "Company3",
				Description: "provide fantastic information.",
				Thumbnail:   "test.com/api.awesome",
				SwaggerURL:  "example.com/api/awesome",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []managementapi.Product{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.req)
			if err != nil {
				t.Errorf("create request body error: %v", err)
			}
			body := bytes.NewReader(bodyBytes)

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/product", body)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			managementapi.PostProduct(w, r)

			rw := w.Result()
			if rw.StatusCode != tt.wantHttpStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantHttpStatus)
			}

			rows, err := db.Queryx("SELECT * from apiinfo WHERE source=$1", tt.req.Source)
			if err != nil {
				t.Errorf("db get products error: %v", err)
				return
			}

			list := []managementapi.Product{}
			for rows.Next() {
				var row managementapi.Product

				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}

				list = append(list, row)
			}

			if diff := cmp.Diff(tt.wantRecords, list, cmpopts.IgnoreFields(managementapi.Product{}, "ID")); diff != "" {
				t.Errorf("db get products responce differs:\n %v", diff)
			}
		})
	}

	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

}
