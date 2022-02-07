package managementapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

var (
	schemaDecoder *schema.Decoder

	UnmarshalJsonErr        = errors.New("failed to parse body as json")
	OtherInputValidationErr = errors.New("input validation failed")
)

const (
	ResultLimitDefault = 50
)

func init() {

	schemaDecoder = schema.NewDecoder()
}

type ResultSet struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type APIInfo struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Source      string `json:"source" db:"source"`
	Description string `json:"description" db:"description"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail"`
	SwaggerURL  string `json:"swagger_url" db:"swagger_url"`
}

type APIInfoList struct {
	List []APIInfo `json:"api_info_list"`
}

type PostAPIInfoReq struct {
	Name        string `json:"name" db:"name" validate:"required"`
	Source      string `json:"source" db:"source" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail" validate:"required,url"`
	SwaggerURL  string `json:"swagger_url" db:"swagger_url" validate:"required,url"`
}

func (pi *PostAPIInfoReq) UnmarshalJSON(data []byte) error {
	type Alias PostAPIInfoReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pi),
	}

	r := bytes.NewReader(data)
	if err := json.NewDecoder(r).Decode(target); err != nil {
		return UnmarshalJsonErr
	}

	if err := ValidateStruct(pi); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			return ve
		} else {
			// unreachable, because ValidateStruct returns ValidationErrors or nil
			return OtherInputValidationErr
		}
	}
	return nil
}

type SearchAPIInfoReq struct {
	Q            string `json:"q" schema:"name" validate:"required,url_encoded"`
	TargetFields string `json:"target_fields" schema:"target_fields"`
	PatternMatch string `json:"pattern_match" schema:"pattern_match"`
	Limit        int    `json:"limit" schema:"limit"`
	Offset       int    `json:"offset" schema:"offset"`
}

func (sr SearchAPIInfoReq) CreateParams() (*SearchAPIInfoParams, error) {
	var err error
	if err = ValidateStruct(sr); err != nil {
		return nil, err
	}
	qSplit := strings.Split(sr.Q, ".")
	for i, v := range qSplit {
		if qSplit[i], err = url.QueryUnescape(v); err != nil {
			return nil, fmt.Errorf("decode string %s error: %w", v, err)
		}
	}

	targetSplit := strings.Split(sr.TargetFields, ".")
	if sr.TargetFields == "" {
		targetSplit = []string{"all"}
	}
	targetFieldExpand := targetSplit
	for _, v := range targetSplit {
		if v == "all" {
			targetFieldExpand = []string{"name", "source", "description"}
			break
		}
	}

	patternMatch := sr.PatternMatch
	if patternMatch == "" {
		patternMatch = "partial"
	}

	limit := sr.Limit
	if limit == 0 {
		limit = ResultLimitDefault
	}

	params := SearchAPIInfoParams{
		Q:            qSplit,
		TargetFields: targetFieldExpand,
		PatternMatch: patternMatch,
		Limit:        limit,
		Offset:       sr.Offset,
	}

	if err = ValidateStruct(params); err != nil {
		return nil, err
	}

	return &params, nil
}

type SearchAPIInfoResult struct {
	APIInfo
	Count int `db:"count"`
}

type SearchAPIInfoMetaData struct {
	ResultSet ResultSet `json:"result_set"`
}

type SearchAPIInfoResp struct {
	APIList               []APIInfo             `json:"api_info_list"`
	SearchAPIInfoMetaData SearchAPIInfoMetaData `json:"metadata"`
}

type SearchAPIInfoParams struct {
	Q            []string `json:"q" validate:"gte=1,dive,ne="`
	TargetFields []string `json:"target_fields" validate:"dive,eq=all|eq=name|eq=description|eq=source"`
	PatternMatch string   `json:"pattern_match" validate:"eq=exact|eq=partial"`
	Limit        int      `json:"limit" validate:"gte=1,lte=100"`
	Offset       int      `json:"offset" validate:"gte=0"`
}

type PostUserReq struct {
	AccountID    string `json:"account_id" db:"account_id" validate:"required,printascii"`
	EmailAddress string `json:"email_address" db:"email_address" validate:"required,email"`
	Password     string `json:"password" db:"password" validate:"required,printascii"`
	Name         string `json:"name" db:"name"`
}

func (pu *PostUserReq) UnmarshalJSON(data []byte) error {
	type Alias PostUserReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pu),
	}

	r := bytes.NewReader(data)
	if err := json.NewDecoder(r).Decode(target); err != nil {
		return UnmarshalJsonErr
	}

	if err := ValidateStruct(pu); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			return ve
		} else {
			// unreachable, because ValidateStruct returns ValidationErrors or nil
			return OtherInputValidationErr
		}
	}
	return nil
}

type User struct {
	ID                string `json:"id" db:"id"`
	AccountID         string `json:"account_id" db:"account_id"`
	EmailAddress      string `json:"email_address" db:"email_address"`
	LoginPasswordHash string `json:"login_password_hash" db:"login_password_hash"`
	Name              string `json:"name" db:"name"`
	PermissionFlag    string `json:"permission_flag" db:"permission_flag"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	UpdatedAt         string `json:"updated_at" db:"updated_at"`
}

/////////////
// product //
/////////////

type Product struct {
	ID              int    `json:"id" db:"id"`
	Name            string `json:"name" db:"name"`
	DisplayName     string `json:"display_name" db:"display_name"`
	Source          string `json:"source" db:"source"`
	Description     string `json:"description" db:"description"`
	Thumbnail       string `json:"thumbnail" db:"thumbnail"`
	IsAvailableCode int    `json:"is_available" db:"is_available"`
}

type PostProductReq struct {
	Name            string       `json:"name" db:"name" validate:"required,alphanum"`
	DisplayName     string       `json:"display_name" db:"display_name"`
	Source          string       `json:"source" db:"source" validate:"required"`
	Description     string       `json:"description" db:"description" validate:"required"`
	Thumbnail       string       `json:"thumbnail" db:"thumbnail" validate:"required,url"`
	Contents        []APIContent `json:"api_contents" validate:"dive"`
	IsAvailable     bool         `json:"is_available"`
	IsAvailableCode int          `json:"-" db:"is_available"`
}

func (pp *PostProductReq) UnmarshalJSON(data []byte) error {
	type Alias PostProductReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pp),
	}

	r := bytes.NewReader(data)
	if err := json.NewDecoder(r).Decode(target); err != nil {
		return UnmarshalJsonErr
	}

	if err := ValidateStruct(pp); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			return ve
		} else {
			// unreachable, because ValidateStruct returns ValidationErrors or nil
			return OtherInputValidationErr
		}
	}
	return nil
}

func (pp PostProductReq) convert() PostProductReq {
	if pp.IsAvailable {
		pp.IsAvailableCode = 1
	} else {
		pp.IsAvailableCode = 0
	}
	return pp
}

type APIContent struct {
	ID          int    `json:"id" db:"id" validate:"required"`
	Description string `json:"description" db:"description"`
}

/////////////
// routing //
/////////////

type PostAPIRoutingReq struct {
	ApiKey     string `json:"api_key" validate:"required"`
	Path       string `json:"path" validate:"required"`
	ForwardURL string `json:"forward_url" validate:"required,url"`
}

func (pr *PostAPIRoutingReq) UnmarshalJSON(data []byte) error {
	type Alias PostAPIRoutingReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pr),
	}

	r := bytes.NewReader(data)
	if err := json.NewDecoder(r).Decode(target); err != nil {
		return UnmarshalJsonErr
	}

	if err := ValidateStruct(pr); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			return ve
		} else {
			// unreachable, because ValidateStruct returns ValidationErrors or nil
			return OtherInputValidationErr
		}
	}
	return nil
}
