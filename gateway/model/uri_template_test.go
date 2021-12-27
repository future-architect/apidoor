package model_test

import (
	"github.com/future-architect/apidoor/gateway/model"
	"testing"
)

func TestURITemplate(t *testing.T) {
	cases := []struct {
		reqPath  string
		template string
		params   []string
		isMatch  bool
	}{
		{
			reqPath:  "/a/b/c",
			template: "/a/b/c",
			params:   []string{},
			isMatch:  true,
		},
		{
			reqPath:  "/a/b/c",
			template: "/a/b/d",
			params:   []string{},
			isMatch:  false,
		},
		{
			reqPath:  "/a/b/c",
			template: "/a/b",
			params:   []string{},
			isMatch:  false,
		},
		{
			reqPath:  "/a/b/c",
			template: "/a/b/c/d",
			params:   []string{},
			isMatch:  false,
		},
		{
			reqPath:  "/a/b/c",
			template: "/a/{test}/c/d",
			params:   []string{},
			isMatch:  false,
		},
		{
			reqPath:  "/a/b/c",
			template: "/a/b/{test}",
			params:   []string{"c"},
			isMatch:  true,
		},
		{
			reqPath:  "/a/b/c",
			template: "/a/{test1}/{test2}",
			params:   []string{"b", "c"},
			isMatch:  true,
		},
	}

	for i, tt := range cases {
		u := model.NewURITemplate(tt.reqPath)
		v := model.NewURITemplate(tt.template)
		params, isMatch := u.Match(v)
		if isMatch != tt.isMatch {
			t.Fatalf("case %d: whether template and request are same or not is wrong", i)
		}
		if params == nil {
			continue
		}
		if err := v.AllocateParameter(params); err != nil {
			t.Fatalf("case %d: get error %s", i, err)
		}
		if v.JoinPath() != tt.reqPath[1:] {
			t.Fatalf("case %d: unexpected path %s, expected %s", i, v.JoinPath(), tt.reqPath[1:])
		}
	}
}
