package managementapi_test

import (
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/go-playground/validator.v8"
)

func TestSearchProductsReq_CreateParams(t *testing.T) {
	tests := []struct {
		name  string
		input managementapi.SearchProductsReq
		want  *managementapi.SearchProductsParams
		// wantErr は期待されるerrorでvalidator.FieldErrorsが返る。validator.FieldErrorは出力に関わるFieldとTagのみ比較する
		wantErr error
	}{
		{
			name: "パーセントエンコードされたクエリを分割して、デコードできる",
			input: managementapi.SearchProductsReq{
				Q:            "a.bc%2ed.efg",
				TargetFields: "name.description",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: &managementapi.SearchProductsParams{
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
			input: managementapi.SearchProductsReq{
				Q:            "abc",
				TargetFields: "all",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: &managementapi.SearchProductsParams{
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
			input: managementapi.SearchProductsReq{
				Q:            "abc",
				TargetFields: "all.name",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: &managementapi.SearchProductsParams{
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
			input: managementapi.SearchProductsReq{
				Q: "abc",
			},
			want: &managementapi.SearchProductsParams{
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
			input: managementapi.SearchProductsReq{
				TargetFields: "name.description",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				"SearchProductsParams.Q[0]": &validator.FieldError{
					Field: "Q[0]",
					Tag:   "ne",
				},
			},
		},
		{
			name: "TargetFieldsに不正な値が含まれている",
			input: managementapi.SearchProductsReq{
				Q:            "abc",
				TargetFields: "name.wrong",
				PatternMatch: "exact",
				Limit:        50,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				"SearchProductsParams.TargetFields[1]": &validator.FieldError{
					Field: "TargetFields[1]",
					Tag:   "eq|eq|eq|eq",
				},
			},
		},
		{
			name: "PatternMatchが不正な値",
			input: managementapi.SearchProductsReq{
				Q:            "abc",
				TargetFields: "all",
				PatternMatch: "wrong",
				Limit:        50,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				"SearchProductsParams.PatternMatch": &validator.FieldError{
					Field: "PatternMatch",
					Tag:   "eq|eq",
				},
			},
		},
		{
			name: "Limitが上限を超えている",
			input: managementapi.SearchProductsReq{
				Q:            "abc",
				TargetFields: "name",
				PatternMatch: "exact",
				Limit:        101,
				Offset:       0,
			},
			want: nil,
			wantErr: validator.ValidationErrors{
				"SearchProductsParams.Limit": &validator.FieldError{
					Field: "Limit",
					Tag:   "lte",
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
			if diff := cmp.Diff(tt.wantErr, err,
				cmpopts.IgnoreFields(validator.FieldError{}, "FieldNamespace", "NameNamespace",
					"Name", "Kind", "ActualTag", "Type", "Param", "Value")); diff != "" {
				t.Errorf("returned error differ:\n%s", diff)
			}
		})
	}
}
