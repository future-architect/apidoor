package usecase

import "github.com/future-architect/apidoor/managementapi/validator"

type ValidateError validator.ValidationErrors

type ClientError struct {
	error
}

func NewClientError(err error) ClientError {
	return ClientError{err}
}

type ServerError struct {
	error
}

func NewServerError(err error) ServerError {
	return ServerError{err}
}
