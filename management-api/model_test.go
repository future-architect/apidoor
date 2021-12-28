package managementapi_test

import (
	"testing"

	"github.com/future-architect/apidoor/managementapi"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

type validateErrorInfo struct {
	field string
	tag   string
}

func TestSearchProductsReq_CreateParams(t *testing.T) {
	tests := []struct {
		name  string
		input managementapi.SearchProductsReq
		want  *managementapi.SearchProductsParams
		// wantErr は期待されるerrorでvalidator.FieldErrorsが返る。validator.FieldErrorは出力に関わるFieldとTagのみ比較する
		wantErr []validateErrorInfo
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
			wantErr: []validateErrorInfo{
				{
					field: "Q[0]",
					tag:   "ne",
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
			wantErr: []validateErrorInfo{
				{
					field: "TargetFields[1]",
					tag:   "eq=all|eq=name|eq=description|eq=source",
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
			wantErr: []validateErrorInfo{
				{
					field: "PatternMatch",
					tag:   "eq=exact|eq=partial",
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
			wantErr: []validateErrorInfo{
				{
					field: "Limit",
					tag:   "lte",
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
			testValidateErrorMatch(t, tt.wantErr, err)
		})
	}
}

func testValidateErrorMatch(t *testing.T, want []validateErrorInfo, got error) {
	gotErr, ok := got.(validator.ValidationErrors)
	if !ok {
		t.Errorf("return error is not validator.ValidationErrors, got: %v", got)
		return
	}

	if len(want) != len(gotErr) {
		t.Errorf("the number of errors is not equal: want=%d, got=%d", len(want), len(gotErr))
		return
	}

	for i := range want {
		if want[i].field != gotErr[i].Field() {
			t.Errorf("field name is not match: want=%s, got=%s", want[i].field, gotErr[i].Field())
		}
		if want[i].tag != gotErr[i].Tag() {
			t.Errorf("tag name is not match: want=%s, got=%s", want[i].tag, gotErr[i].Tag())
		}
	}

}
