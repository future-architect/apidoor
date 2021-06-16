package gateway

import (
	"errors"
	"strings"
)

type block struct {
	value   string
	isparam bool
}

type URITemplate struct {
	path []block
}

func NewURITemplate(path string) *URITemplate {
	u := new(URITemplate)
	slice := strings.Split(path[1:], "/")
	for _, v := range slice {
		u.path = append(u.path, block{
			value:   v,
			isparam: strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}"),
		})
	}

	return u
}

/*
func (u *URITemplate) Init(path string) {
	slice := strings.Split(path[1:], "/")
	for _, v := range slice {
		u.path = append(u.path, block{
			value:   v,
			isparam: strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}"),
		})
	}
}
*/

func (u *URITemplate) TemplateMatch(t URITemplate) (bool, []string) {
	var params []string
	if len(u.path) != len(t.path) {
		return false, nil
	}

	for i := 0; i < len(u.path); i++ {
		if t.path[i].isparam {
			params = append(params, u.path[i].value)
		} else if u.path[i].value != t.path[i].value {
			return false, nil
		}
	}

	return true, params
}

func (u *URITemplate) JoinPath() string {
	var s []string
	for _, v := range u.path {
		s = append(s, v.value)
	}

	return strings.Join(s, "/")
}

func (u *URITemplate) AssignParameter(s []string) error {
	var indices []int
	for i, v := range u.path {
		if v.isparam {
			indices = append(indices, i)
		}
	}

	if len(indices) != len(s) {
		return errors.New("number of parameters doesn't match")
	}

	for i, v := range indices {
		u.path[v].value = s[i]
		u.path[v].isparam = false
	}

	return nil
}
