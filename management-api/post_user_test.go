package managementapi_test

import (
	"bytes"
	"encoding/json"
	"github.com/future-architect/apidoor/managementapi"
	"github.com/future-architect/apidoor/managementapi/model"
	"github.com/future-architect/apidoor/managementapi/validator"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestPostUser(t *testing.T) {
	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}

	hashRegex := regexp.MustCompile(`\$2a\$\w+\$[ -~]+`)

	tests := []struct {
		name               string
		contentType        string
		req                model.PostUserReq
		wantHttpStatus     int
		wantBadRequestResp *validator.BadRequestResp
		//wantRecord は期待されるDB作成レコードの値、idは比較対象外
		wantRecords []model.User
	}{
		{
			name:        "正常に登録できる",
			contentType: "application/json",
			req: model.PostUserReq{
				AccountID:    "user",
				EmailAddress: "test00@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus:     http.StatusCreated,
			wantBadRequestResp: nil,
			wantRecords: []model.User{
				{
					AccountID:      "user",
					EmailAddress:   "test00@example.com",
					Name:           "full name",
					PermissionFlag: "00",
				},
			},
		},
		{
			name:        "パスワードに記号が含まれており、正常に登録できる",
			contentType: "application/json",
			req: model.PostUserReq{
				AccountID:    "user1",
				EmailAddress: "test01@example.com",
				Password:     "p@ss12Word",
				Name:         "full name",
			},
			wantHttpStatus:     http.StatusCreated,
			wantBadRequestResp: nil,
			wantRecords: []model.User{
				{
					AccountID:      "user1",
					EmailAddress:   "test01@example.com",
					Name:           "full name",
					PermissionFlag: "00",
				},
			},
		},
		{
			name:        "名前が設定されていなくても、正常に登録できる",
			contentType: "application/json",
			req: model.PostUserReq{
				AccountID:    "user2",
				EmailAddress: "test02@example.com",
				Password:     "password",
				Name:         "",
			},
			wantHttpStatus:     http.StatusCreated,
			wantBadRequestResp: nil,
			wantRecords: []model.User{
				{
					AccountID:      "user2",
					EmailAddress:   "test02@example.com",
					Name:           "",
					PermissionFlag: "00",
				},
			},
		},
		{
			name:        "必須項目に空欄があるとき、登録できない",
			contentType: "application/json",
			req: model.PostUserReq{
				AccountID:    "",
				EmailAddress: "test03@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantBadRequestResp: &validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "account_id",
						ConstraintType: "required",
						Message:        "required field, but got empty",
						Got:            "",
					},
				},
			},
			wantRecords: []model.User{},
		},
		{
			name:        "account_idにprintable ascii以外の文字が含まれていたとき、登録できない",
			contentType: "application/json",
			req: model.PostUserReq{
				AccountID:    "userユーザー",
				EmailAddress: "test04@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantBadRequestResp: &validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "account_id",
						ConstraintType: "printascii",
						Message:        "input value, userユーザー, does not satisfy the format, printascii",
						Got:            "userユーザー",
					},
				},
			},
			wantRecords: []model.User{},
		},
		{
			name:        "email_addressの文字列がメールアドレスとして不正であるとき、登録できない",
			contentType: "application/json",
			req: model.PostUserReq{
				AccountID:    "user",
				EmailAddress: "test05.@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantBadRequestResp: &validator.BadRequestResp{
				Message: "input validation error",
				ValidationErrors: &validator.ValidationErrors{
					{
						Field:          "email_address",
						ConstraintType: "email",
						Message:        "input value, test05.@example.com, does not satisfy the format, email",
						Got:            "test05.@example.com",
					},
				},
			},
			wantRecords: []model.User{},
		},
		{
			name:        "Content-Typeがapplication/json以外であるとき、登録できない",
			contentType: "text/plain",
			req: model.PostUserReq{
				AccountID:    "user",
				EmailAddress: "test06@example.com",
				Password:     "password",
				Name:         "full name",
			},
			wantHttpStatus: http.StatusBadRequest,
			wantBadRequestResp: &validator.BadRequestResp{
				Message: `unexpected request Content-Type, it must be "application/json"`,
			},
			wantRecords: []model.User{},
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

			r := httptest.NewRequest(http.MethodPost, "localhost:3000/mgmt/users", body)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			managementapi.PostUser(w, r)

			rw := w.Result()
			defer rw.Body.Close()

			if rw.StatusCode != tt.wantHttpStatus {
				t.Errorf("wrong http status code: got %d, want %d", rw.StatusCode, tt.wantHttpStatus)
			}

			resp, err := io.ReadAll(rw.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rw.StatusCode == http.StatusBadRequest {
				testBadRequestResp(t, tt.wantBadRequestResp, resp)
				return
			}
			if rw.StatusCode != http.StatusCreated {
				return
			}

			rows, err := db.Queryx("SELECT * from apiuser WHERE email_address=$1", tt.req.EmailAddress)
			if err != nil {
				t.Errorf("db get users error: %v", err)
				return
			}

			list := []model.User{}
			for rows.Next() {
				var row model.User

				if err := rows.StructScan(&row); err != nil {
					t.Errorf("reading row error: %v", err)
					return
				}

				list = append(list, row)
			}

			if diff := cmp.Diff(tt.wantRecords, list,
				cmpopts.IgnoreFields(model.User{}, "ID", "LoginPasswordHash",
					"CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf("db get users responce differs:\n %v", diff)
			}

			// checking that passwords are stored in a hash
			for _, v := range list {
				if !hashRegex.Match([]byte(v.LoginPasswordHash)) {
					t.Errorf("password hash format is wrong, got: %s", v.LoginPasswordHash)
				}
			}

		})
	}

	if _, err := db.Exec("DELETE FROM apiuser"); err != nil {
		t.Fatal(err)
	}

}

func testBadRequestResp(t *testing.T, want *validator.BadRequestResp, got []byte) {
	var gotBody validator.BadRequestResp
	if err := json.Unmarshal(got, &gotBody); err != nil {
		t.Errorf("parsing body as BadRequestResp failed: %v", err)
		return
	}

	if diff := cmp.Diff(gotBody, *want); diff != "" {
		t.Errorf("bad request response differs:\n%v", diff)
	}

}
