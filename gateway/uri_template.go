package gateway

import (
	"errors"
	"path"
	"strings"
)

type block struct {
	value   string
	isParam bool
}

type URITemplate struct {
	path []block
}

func NewURITemplate(path string) *URITemplate {
	if len(path) == 0 {
		panic("empty path")
	} else if len(path) == 1 {
		return &URITemplate{}
	}

	u := &URITemplate{}
	slice := strings.Split(strings.Trim(path, "/"), "/")
	for _, v := range slice {
		isParam := strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")
		var value string
		if isParam {
			value = v[1 : len(v)-1]
		} else {
			value = v
		}
		u.path = append(u.path, block{
			value:   value,
			isParam: isParam,
		})
	}

	return u
}

func (u *URITemplate) Match(t URITemplate) (map[string]string, bool) {
	if len(u.path) != len(t.path) {
		return nil, false
	}

	params := make(map[string]string, len(u.path))

	for i := 0; i < len(u.path); i++ {
		if t.path[i].isParam {
			params[t.path[i].value] = u.path[i].value
		} else if u.path[i].value != t.path[i].value {
			return nil, false
		}
	}

	return params, true
}

func (u *URITemplate) JoinPath() string {
	s := make([]string, 0, len(u.path))
	for _, v := range u.path {
		s = append(s, v.value)
	}

	return path.Join(s...)
}

func (u *URITemplate) AllocateParameter(m map[string]string) error {
	for i, block := range u.path {
		if block.isParam {
			if v, ok := m[block.value]; !ok {
				return errors.New("no such parameter")
			} else {
				u.path[i].value = v
			}
		}
	}

	return nil
}
