package managementapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPostAPIInfo(t *testing.T) {
	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		contentType    string
		req            model.PostAPIInfoReq
		wantHttpStatus int
		//wantRecord は期待されるDB作成レコードの値、idは比較対象外
		wantRecords []model.APIInfo
		wantResp    interface{}
	}{
		{
			name:        "api infoを登録できる",
			contentType: "application/json",
			req: model.PostAPIInfoReq{
				Name:        "Awesome API",
				Source:      "Company1",
				Description: "provide fantastic information.",
				Thumbnail:   "http://example.com/api.awesome",
				SwaggerURL:  "http://example.com/api/awesome",
			},
			wantHttpStatus: http.StatusCreated,
			wantRecords: []model.APIInfo{
				{
					Name:        "Awesome API",
					Source:      "Company1",
					Description: "provide fantastic information.",
					Thumbnail:   "http://example.com/api.awesome",
					SwaggerURL:  "http://example.com/api/awesome",
				},
			},
			wantResp: "Created",
		},
		{
			name:        "Fieldに空文字列がある場合は登録できない",
			contentType: "application/json",
			req: model.PostAPIInfoReq{
				Name:        "",
				Source:      "Company2",
				Description: "provide fantastic information.",
				Thumbnail:   "http://example.com/api.awesome",
				SwaggerURL:  "http://example.com/api/awesome",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []model.APIInfo{},
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
			req: model.PostAPIInfoReq{
				Name:        "wrong content-type",
				Source:      "Company3",
				Description: "provide fantastic information.",
				Thumbnail:   "http://example.com/api.awesome",
				SwaggerURL:  "http://example.com/api/awesome",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantRecords:    []model.APIInfo{},
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
			managementapi.PostAPIInfo(w, r)

			rw := w.Result()

			resp, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rw.StatusCode != tt.wantHttpStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantHttpStatus)
			}

			rows, err := db.Queryx("SELECT * from apiinfo WHERE source=$1", tt.req.Source)
			if err != nil {
				t.Errorf("db get api info error: %v", err)
				return
			}

			list := []model.APIInfo{}
			for rows.Next() {
				var row model.APIInfo

				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}

				list = append(list, row)
			}

			if diff := cmp.Diff(tt.wantRecords, list, cmpopts.IgnoreFields(model.APIInfo{}, "ID")); diff != "" {
				t.Errorf("db get list of api info responce differs:\n %v", diff)
			}
			switch tt.wantResp.(type) {
			case string:
				if tt.wantResp != string(resp) {
					t.Errorf("response body is not %s, got %s", tt.wantResp, string(resp))
				}
			case validator.BadRequestResp:
				want := tt.wantResp.(validator.BadRequestResp)
				testBadRequestResp(t, &want, resp)
			default:
				t.Errorf("type of wantResp is not unsupported")
			}
		})
	}

	if _, err := db.Exec("DELETE FROM apiinfo"); err != nil {
		t.Fatal(err)
	}

}
