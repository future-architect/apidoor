package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/validator"
)

// TODO: ../model.goの移動

var (
	UnmarshalJsonErr        = errors.New("failed to parse body as json")
	OtherInputValidationErr = errors.New("input validation failed")
)

////////////////
// api tokens //
////////////////

type ParamType string

const (
	Header          ParamType = "header"
	Query                     = "query"
	BodyFormEncoded           = "body_form_encoded"
)

type AccessToken struct {
	ParamType ParamType `dynamo:"param_type" json:"param_type" validate:"required,eq=header|eq=query|eq=body_from_encoded"`
	Key       string    `dynamo:"key" json:"key" validate:"required"`
	Value     string    `dynamo:"value" json:"value" validate:"required"`
}

type PostAPITokenReq struct {
	APIKey       string        `json:"api_key" validate:"required"`
	Path         string        `json:"path" validate:"required"`
	AccessTokens []AccessToken `json:"tokens" validate:"required,dive,required"`
}

func (pp *PostAPITokenReq) UnmarshalJSON(data []byte) error {
	type Alias PostAPITokenReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pp),
	}

	r := bytes.NewReader(data)
	if err := json.NewDecoder(r).Decode(target); err != nil {
		return fmt.Errorf("api token req: %s %w", err.Error(), UnmarshalJsonErr)
	}

	if err := validator.ValidateStruct(pp); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return ve
		} else {
			// unreachable, because ValidateStruct returns ValidationErrors or nil
			return OtherInputValidationErr
		}
	}
	return nil
}
