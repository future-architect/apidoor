package model

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnauthorizedRequest = &MyError{Message: "unauthorized request"}

type MyError struct {
	Message string `json:"message"`
}

func (err *MyError) Error() string {
	return fmt.Sprintf("error: %s", err.Message)
}

type Field struct {
	// Template is a gateway path
	Template URITemplate
	// Path is a  destination api path
	Path          URITemplate
	ForwardSchema string
	// Num represents the recent number of api calls.
	Num int
	// Max represents the maximum api call limit if Max's type is int.
	// If it is "-", the number of the api calls has no limit.
	Max interface{}
}

func (f Field) createForwardURL(query map[string]string) string {
	nodes := make([]string, 0, len(f.Path.path))
	for _, v := range f.Path.path {
		if v.isParam {
			nodes = append(nodes, query[v.value])
		} else {
			nodes = append(nodes, v.value)
		}
	}
	var schema string
	if f.ForwardSchema != "" {
		schema = f.ForwardSchema + "://"
	}
	return schema + strings.Join(nodes, "/")
}

type Fields []Field

type FieldResult struct {
	Field        Field
	ForwardURL   string
	TemplatePath string
}

func (f Fields) LookupTemplate(path string) (*FieldResult, error) {
	u := NewURITemplate(path)
	for _, v := range f {
		if params, ok := u.Match(v.Template); ok {
			forwardURL := v.createForwardURL(params)
			return &FieldResult{Field: v, ForwardURL: forwardURL, TemplatePath: v.Template.JoinPath()}, nil
		}
	}
	return nil, ErrUnauthorizedRequest // Not found path
}

func (f Fields) CheckAPILimit(path string) error {
	for _, field := range f {
		if field.Template.JoinPath() == path {
			switch max := field.Max.(type) {
			case int:
				if field.Num >= max {
					return errors.New("limit exceeded")
				}
			case string:
				if max != "-" {
					return errors.New("unexpected limit value")
				}
			default:
				return errors.New("unexpected limit value")
			}

			return nil
		}
	}

	return nil
}
