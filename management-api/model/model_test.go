package model

import (
	"github.com/future-architect/apidoor/managementapi/validator"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSearchProductReq_CreateParams(t *testing.T) {
	tests := []struct {
		name  string
		input SearchProductReq
		want  *SearchProductParams
		// wantErr は期待されるerrorでvalidator.FieldErrorsが返る。validator.FieldErrorは出力に関わるFieldとTagのみ比較する
		wantErr validator.ValidationErrors
	}{
		{
			name: "パーセントエンコードされたクエリを分割して、デコードできる",
			input: SearchProductReq{
				Q:            "a.bc%2ed.efg",
				TargetFields: "name.description",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: &SearchProductParams{
				Q:            []string{"a", "bc.d", "efg"},
				TargetFields: []string{"name", "description"},
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			wantErr: nil,
		},
		{
			name: "target_fieldsがallのとき、展開される",
			input: SearchProductReq{
				Q:            "abc",
				TargetFields: "all",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: &SearchProductParams{
				Q:            []string{"abc"},
				TargetFields: []string{"name", "source", "description"},
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			wantErr: nil,
		},
		{
			name: "target_fieldsがallを含めて複数あるとき、all単独のときと同様に展開される",
			input: SearchProductReq{
				Q:            "abc",
				TargetFields: "all.name",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: &SearchProductParams{
				Q:            []string{"abc"},
				TargetFields: []string{"name", "source", "description"},
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			wantErr: nil,
		},
		{
			name: "未指定時に既定の値が設定される",
			input: SearchProductReq{
				Q: "abc",
			},
			want: &SearchProductParams{
				Q:            []string{"abc"},
				TargetFields: []string{"name", "source", "description"},
				PatternMatch: "partial",
				Limit:        50,
				Offset:       0,
			},
			wantErr: nil,
		},
		{
			name: "クエリが空",
			input: SearchProductReq{
				TargetFields: "name.description",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				{
					Field:          "q",
					ConstraintType: "required",
					Message:        "required field, but got empty",
					Got:            "",
				},
			},
		},
		{
			name: "クエリに検索語に空文字列がある",
			input: SearchProductReq{
				Q: "abc.",
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				{
					Field:          "q[1]",
					ConstraintType: "ne",
					Message:        "input value is , but it must be not equal to ",
					Got:            "",
				},
			},
		},
		{
			name: "クエリがURL encodedとして不正",
			input: SearchProductReq{
				Q: "a%g3bc",
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				{
					Field:          "q",
					ConstraintType: "url_encoded",
					Message:        "input value, a%g3bc, does not satisfy the format, url_encoded",
					Got:            "a%g3bc",
				},
			},
		},
		{
			name: "TargetFieldsに不正な値が含まれている",
			input: SearchProductReq{
				Q:            "abc",
				TargetFields: "name.wrong",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				{
					Field:          "target_fields[1]",
					ConstraintType: "enum",
					Message:        "input value is wrong, but it must be one of the following values: [all name description source]",
					Enum:           []string{"all", "name", "description", "source"},
					Got:            "wrong",
				},
			},
		},
		{
			name: "PatternMatchが不正な値",
			input: SearchProductReq{
				Q:            "abc",
				TargetFields: "all",
				PatternMatch: "wrong",
				Limit:        50,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				{
					Field:          "pattern_match",
					ConstraintType: "enum",
					Message:        "input value is wrong, but it must be one of the following values: [exact partial]",
					Enum:           []string{"exact", "partial"},
					Got:            "wrong",
				},
			},
		},
		{
			name: "Limitが上限を超えている",
			input: SearchProductReq{
				Q:            "abc",
				TargetFields: "name",
				PatternMatch: "exact",
				Limit:        101,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				{
					Field:          "limit",
					ConstraintType: "lte",
					Message:        "input value is 101, but it must be less than or equal to 100",
					Lte:            "100",
					Got:            101,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.input.CreateParams()
			if diff := cmp.Diff(tt.want, resp); diff != "" {
				t.Errorf("retruned struct differ:\n%s", diff)
			}

			if err == nil {
				if tt.wantErr != nil {
					t.Errorf("returned error is nil, but expected error is not nil: %v", tt.wantErr)
				}
				return
			}
			testValidateErrors(t, tt.wantErr, err)
		})
	}
}

func testValidateErrors(t *testing.T, want validator.ValidationErrors, got error) {
	gotErr, ok := got.(validator.ValidationErrors)
	if !ok {
		t.Errorf("return error is not validator.ValidationErrors, got: %v", got)
		return
	}

	if diff := cmp.Diff(want, gotErr); diff != "" {
		t.Errorf("ValidationErrors differ:\n%v", diff)
	}

}
