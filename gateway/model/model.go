package model

import (
	"errors"
	"fmt"
)

var ErrUnauthorizedRequest = &MyError{Message: "unauthorized request"}

type MyError struct {
	Message string `json:"message"`
}

func (err *MyError) Error() string {
	return fmt.Sprintf("error: %s", err.Message)
}

type Field struct {
	Template      URITemplate
	Path          URITemplate
	ForwardSchema string
	Num           int
	Max           interface{}
}

type Fields []Field

func (f Fields) URI(path string) (string, error) {
	u := NewURITemplate(path)
	for _, v := range f {
		if _, ok := u.Match(v.Template); ok {
			if v.ForwardSchema == "" {
				return v.Path.JoinPath(), nil
			}
			return v.ForwardSchema + "://" + v.Path.JoinPath(), nil
		}
	}
	return "", ErrUnauthorizedRequest // Not found path
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
