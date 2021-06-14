package gateway

import "strings"

type block struct {
	value   string
	isparam bool
}

type URITemplate struct {
	path []block
}

func (u *URITemplate) Init(path string) {
	slice := strings.Split(path, "/")
	for _, v := range slice {
		u.path = append(u.path, block{
			value:   v,
			isparam: strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}"),
		})
	}
}

func (u *URITemplate) TemplateMatch(t URITemplate) (bool, []string) {
	var params []string
	if len(u.path) != len(t.path) {
		return false, []string{}
	}

	for i := 0; i < len(u.path); i++ {
		if t.path[i].isparam {
			params = append(params, u.path[i].value)
		} else if u.path[i].value != t.path[i].value {
			return false, []string{}
		}
	}

	return true, params
}
