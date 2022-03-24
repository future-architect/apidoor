package swaggerparser

import (
	"errors"
	"fmt"
)

type ErrorType string

// TODO: 名前の検討(他とはclientとserverの意味が逆)
const (
	FetchClientError ErrorType = "fetch file client error"
	FetchServerError ErrorType = "fetch file server error"
	FileParseError   ErrorType = "parse file error"
	FormatError      ErrorType = "format error"
	OtherError       ErrorType = "other error"
)

type Error struct {
	ErrorType ErrorType
	Message   error
}

func newError(errorType ErrorType, message error) Error {
	return Error{
		ErrorType: errorType,
		Message:   message,
	}
}

func newErrorString(errorType ErrorType, message string) Error {
	return newError(errorType, errors.New(message))
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorType, e.Message.Error())
}
