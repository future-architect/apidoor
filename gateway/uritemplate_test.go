package gateway_test

import (
	"gateway"
	"testing"
)

type templatetest struct {
	reqpath  string
	template string
	params   []string
	ismatch  bool
}

var templatedata = []templatetest{
	{
		reqpath:  "/a/b/c",
		template: "/a/b/c",
		params:   []string{},
		ismatch:  true,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b/d",
		params:   []string{},
		ismatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b",
		params:   []string{},
		ismatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b/c/d",
		params:   []string{},
		ismatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/{test}/c/d",
		params:   []string{},
		ismatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b/{test}",
		params:   []string{"c"},
		ismatch:  true,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/{test1}/{test2}",
		params:   []string{"b", "c"},
		ismatch:  true,
	},
}

func TestURITemplate(t *testing.T) {
	for i, tt := range templatedata {
		var u, v gateway.URITemplate
		u.Init(tt.reqpath)
		v.Init(tt.template)
		ismatch, params := u.TemplateMatch(v)
		if ismatch != tt.ismatch {
			t.Fatalf("case %d: whether template and request are same or not is wrong", i)
		}
		if len(params) != len(tt.params) {
			t.Fatalf("case %d: expected matched params size %d, get %d", i, len(tt.params), len(params))
		}
		for j, value := range params {
			if value != tt.params[j] {
				t.Fatalf("case %d: unexpected matched param %s, expected %s", i, value, tt.params[j])
			}
		}
	}
}
