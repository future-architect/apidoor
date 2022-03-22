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
// products //
//////////////

// TODO: productの管理者情報の追加 https://github.com/future-architect/apidoor/issues/79

type Product struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	// OwnerID         int    `json:"owner" db:"owner"`
	DisplayName     string `json:"display_name" db:"display_name"`
	Source          string `json:"source" db:"source"`
	Description     string `json:"description" db:"description"`
	Thumbnail       string `json:"thumbnail" db:"thumbnail"`
	BasePath        string `json:"base_path" db:"base_path"`
	SwaggerURL      string `json:"swagger_url" db:"swagger_url"`
	IsAvailableCode int    `json:"is_available" db:"is_available"`
	CreatedAt       string `json:"created_at" db:"created_at"`
	UpdatedAt       string `json:"updated_at" db:"updated_at"`
}

type ProductList struct {
	List []Product `json:"product_list"`
}

type PostProductReq struct {
	Name   string `json:"name" db:"name" validate:"required"`
	Source string `json:"source" db:"source" validate:"required"`
	// OwnerID     *int   `json:"owner_id" db:"owner_id" validate:"required"`
	DisplayName string `json:"display_name" db:"display_name" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail" validate:"required,url"`
	SwaggerURL  string `json:"swagger_url" db:"swagger_url" validate:"required,url"`
	IsAvailable bool   `json:"is_available"`
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

func (pp *PostProductReq) DBParam(basePath string) PostProductDB {
	isAvailable := 0
	if pp.IsAvailable {
		isAvailable = 1
	}

	return PostProductDB{
		Name:        pp.Name,
		Source:      pp.Source,
		DisplayName: pp.DisplayName,
		Description: pp.Description,
		Thumbnail:   pp.Thumbnail,
		BasePath:    basePath,
		SwaggerURL:  pp.SwaggerURL,
		IsAvailable: isAvailable,
	}
}

type PostProductDB struct {
	Name   string `db:"name"`
	Source string `db:"source"`
	// OwnerID     int    `db:"owner_id"`
	DisplayName string `db:"display_name"`
	Description string `db:"description"`
	Thumbnail   string `db:"thumbnail"`
	BasePath    string `db:"base_path"`
	SwaggerURL  string `db:"swagger_url"`
	IsAvailable int    `db:"is_available"`
}

type SearchProductReq struct {
	Q            string `json:"q" schema:"name" validate:"required,url_encoded"`
	TargetFields string `json:"target_fields" schema:"target_fields"`
	PatternMatch string `json:"pattern_match" schema:"pattern_match"`
	Limit        int    `json:"limit" schema:"limit"`
	Offset       int    `json:"offset" schema:"offset"`
}

func (sr SearchProductReq) CreateParams() (*SearchProductParams, error) {
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

	params := SearchProductParams{
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

type SearchProductResult struct {
	Product
	Count int `db:"count"`
}

type SearchProductMetaData struct {
	ResultSet ResultSet `json:"result_set"`
}

type SearchProductResp struct {
	ProductList           []Product             `json:"product_list"`
	SearchProductMetaData SearchProductMetaData `json:"metadata"`
}

type SearchProductParams struct {
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
	UserAccountID string `json:"user_account_id" validate:"required,printascii"`
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

type APIKeyContractProductAuthorized struct {
	ApiKeyID          int
	ContractProductID int
	CreatedAt         string `json:"created_at" db:"created_at"`
	UpdatedAt         string `json:"updated_at" db:"updated_at"`
}

type PostAPIKeyProductsReq struct {
	ApiKeyID  *int                         `json:"apikey_id" validate:"required"`
	Contracts []AuthorizedContractProducts `json:"contracts" validate:"required,gte=1,dive"`
}

func (pkp PostAPIKeyProductsReq) ContractIDMap() map[int]AuthorizedContractProducts {
	ret := make(map[int]AuthorizedContractProducts)
	for _, contract := range pkp.Contracts {
		ret[contract.ContractID] = contract
	}
	return ret
}

func (pkp *PostAPIKeyProductsReq) UnmarshalJSON(data []byte) error {
	type Alias PostAPIKeyProductsReq
	target := &struct {
		*Alias
	}{
		Alias: (*Alias)(pkp),
	}
	fmt.Println(string(data))
	return validator.UnmarshalJSON(pkp, data, target)
}

type AuthorizedContractProducts struct {
	ContractID int `json:"contract_id" validate:"required"`

	// ProductIDs is the list of product which is linked to the key
	// if this field is nil or empty, all products the contract contains will be linked
	ProductIDs []int `json:"product_ids,omitempty"`
}

type ContractProductDB struct {
	ID         int `db:"id"`
	ContractID int `db:"contract_id"`
	ProductID  int `db:"product_id"`
}
