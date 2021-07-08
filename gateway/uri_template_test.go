package gateway_test

import (
	"gateway"
	"testing"
)

type templatetest struct {
	reqpath  string
	template string
	params   []string
	isMatch  bool
}

var templatedata = []templatetest{
	{
		reqpath:  "/a/b/c",
		template: "/a/b/c",
		params:   []string{},
		isMatch:  true,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b/d",
		params:   []string{},
		isMatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b",
		params:   []string{},
		isMatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b/c/d",
		params:   []string{},
		isMatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/{test}/c/d",
		params:   []string{},
		isMatch:  false,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/b/{test}",
		params:   []string{"c"},
		isMatch:  true,
	},
	{
		reqpath:  "/a/b/c",
		template: "/a/{test1}/{test2}",
		params:   []string{"b", "c"},
		isMatch:  true,
	},
}

func TestURITemplate(t *testing.T) {
	for i, tt := range templatedata {
		u := gateway.NewURITemplate(tt.reqpath)
		v := gateway.NewURITemplate(tt.template)
		params, isMatch := u.TemplateMatch(*v)
		if isMatch != tt.isMatch {
			t.Fatalf("case %d: whether template and request are same or not is wrong", i)
		}
		if params == nil {
			continue
		}
		if err := v.AllocateParameter(params); err != nil {
			t.Fatalf("case %d: get error %s", i, err)
		}
		if v.JoinPath() != tt.reqpath[1:] {
			t.Fatalf("case %d: unexpected path %s, expected %s", i, v.JoinPath(), tt.reqpath[1:])
		}
	}
}
