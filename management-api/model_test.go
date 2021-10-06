package managementapi_test

import (
	"github.com/google/go-cmp/cmp"
	"log"
	"managementapi"
	"testing"
)

func TestSearchProductsReq_CreateParams(t *testing.T) {
	tests := []struct {
		name  string
		input managementapi.SearchProductsReq
		want  *managementapi.SearchProductsParams
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.input.CreateParams()
			if diff := cmp.Diff(tt.want, resp); diff != "" {
				t.Errorf("retruned struct differ:\n%s", diff)
			}
			if err != nil {
				log.Println(err)
			}
		})
	}
}
