package model

import (
	"errors"
	"fmt"
	"github.com/future-architect/apidoor/managementapi/validator"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

var (
	SchemaDecoder *schema.Decoder

	UnmarshalJsonErr = errors.New("failed to parse body as json")
)

const (
	ResultLimitDefault = 50
)

func init() {

	SchemaDecoder = schema.NewDecoder()
}

////////////
// common //
////////////

type ResultSet struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// EmptyResp is used for swaggo to indicate empty response (204 no content etc.)
type EmptyResp struct{}

//////////////
// api info //
//////////////

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
	return validator.UnmarshalJSON(pi, data, target)
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
	if err = validator.ValidateStruct(sr); err != nil {
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

	if err = validator.ValidateStruct(params); err != nil {
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

//////////////
// api user //
//////////////

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
	return validator.UnmarshalJSON(pu, data, target)
}

type User struct {
	ID                int    `json:"id" db:"id"`
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
	CreatedAt       string `json:"created_at" db:"created_at"`
	UpdatedAt       string `json:"updated_at" db:"updated_at"`
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
	return validator.UnmarshalJSON(pp, data, target)
}

func (pp PostProductReq) Convert() PostProductReq {
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

//////////////
// contract //
//////////////

//TODO: 他にcontractに含めるべきカラムの検討

type Contract struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type PostContractReq struct {
	UserAccountID string              `json:"user_id" validate:"required,printascii"`
	Products      []*ContractProducts `json:"products" validate:"required,gte=1,dive,required"`
}

func (pr *PostContractReq) UnmarshalJSON(data []byte) error {
	type Alias PostContractReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pr),
	}
	return validator.UnmarshalJSON(pr, data, target)
}

type ContractProducts struct {
	ProductName string `json:"product_name" validate:"required"`
	Description string `json:"description"`
}

type PostContractDB struct {
	UserID   int
	Products []*ContractProductsDB
}

type ContractProductsDB struct {
	ProductID   int    `db:"product_name"`
	Description string `db:"description"`
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
	return validator.UnmarshalJSON(pr, data, target)
}

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
	return validator.UnmarshalJSON(pp, data, target)
}

type DeleteAPITokenReq struct {
	APIKey string `json:"api_key" schema:"api_key" validate:"required"`
	Path   string `json:"path" schema:"path" validate:"required"`
}

func (da *DeleteAPITokenReq) UnmarshalJSON(data []byte) error {
	type Alias DeleteAPITokenReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(da),
	}
	return validator.UnmarshalJSON(da, data, target)
}

/////////////
// api key //
/////////////

type APIKey struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	AccessKey string `json:"access_key" db:"access_key"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type PostAPIKeyReq struct {
	UserAccountID string `json:"user_account_id" validator:"required,printascii"`
}

func (pk *PostAPIKeyReq) UnmarshalJSON(data []byte) error {
	type Alias PostAPIKeyReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pk),
	}
	return validator.UnmarshalJSON(pk, data, target)
}

type PostAPIKeyResp struct {
	UserAccountID string `json:"user_account_id"`
	AccessKey     string `json:"access_key"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
