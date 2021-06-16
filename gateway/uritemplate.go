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
		isparam := strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}")
		var value string
		if isparam {
			value = v[1 : len(v)-1]
		} else {
			value = v
		}
		u.path = append(u.path, block{
			value,
			isparam,
		})
	}

	return u
}

func (u *URITemplate) TemplateMatch(t URITemplate) (bool, map[string]string) {
	params := make(map[string]string)
	if len(u.path) != len(t.path) {
		return false, nil
	}

	for i := 0; i < len(u.path); i++ {
		if t.path[i].isparam {
			params[t.path[i].value] = u.path[i].value
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

func (u *URITemplate) AllocateParameter(m map[string]string) error {
	for i, block := range u.path {
		if block.isparam {
			if v, ok := m[block.value]; !ok {
				return errors.New("no such parameter")
			} else {
				u.path[i].value = v
			}
		}
	}

	return nil
}
