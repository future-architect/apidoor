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
		u := gateway.NewURITemplate(tt.reqpath)
		v := gateway.NewURITemplate(tt.template)
		ismatch, params := u.TemplateMatch(*v)
		if ismatch != tt.ismatch {
			t.Fatalf("case %d: whether template and request are same or not is wrong", i)
		}
		if params == nil {
			continue
		}
		if err := v.AssignParameter(params); err != nil {
			t.Fatalf("case %d: get error %s", i, err)
		}
		if v.JoinPath() != tt.reqpath[1:] {
			t.Fatalf("case %d: unexpected path %s, expected %s", i, v.JoinPath(), tt.reqpath[1:])
		}
	}
}
