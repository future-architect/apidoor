package managementapi

import (
	"fmt"
	"gopkg.in/go-playground/validator.v8"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

var (
	validate      *validator.Validate
	schemaDecoder *schema.Decoder
)

const (
	ResultLimitDefault = 50
)

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)

	schemaDecoder = schema.NewDecoder()
}

type ResultSet struct {
	Count  int `json:"count"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Product struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Source      string `json:"source" db:"source"`
	Description string `json:"description" db:"description"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail"`
	SwaggerURL  string `json:"swagger_url" db:"swagger_url"`
}

type Products struct {
	Products []Product `json:"products"`
}

type PostProductReq struct {
	Name        string `json:"name" db:"name" validate:"required"`
	Source      string `json:"source" db:"source" validate:"required"`
	Description string `json:"description" db:"description" validate:"required"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail" validate:"required"`
	SwaggerURL  string `json:"swagger_url" db:"swagger_url" validate:"required"`
}

type SearchProductsReq struct {
	Q            string `schema:"name"`
	TargetFields string `schema:"target_fields"`
	PatternMatch string `schema:"pattern_match"`
	Limit        int    `schema:"limit"`
	Offset       int    `schema:"offset"`
}

func (sr SearchProductsReq) CreateParams() (*SearchProductsParams, error) {
	var err error
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

	params := SearchProductsParams{
		Q:            qSplit,
		TargetFields: targetFieldExpand,
		PatternMatch: patternMatch,
		Limit:        limit,
		Offset:       sr.Offset,
	}

	if err = validate.Struct(params); err != nil {
		return nil, err
	}

	return &params, nil
}

type SearchProductsResult struct {
	Product
	Count int `db:"count"`
}

type SearchProductsMetaData struct {
	ResultSet ResultSet `json:"result_set"`
}

type SearchProductsResp struct {
	Products               []Product              `json:"products"`
	SearchProductsMetaData SearchProductsMetaData `json:"metadata"`
}

type SearchProductsParams struct {
	Q            []string `validate:"gte=1,dive,ne="`
	TargetFields []string `validate:"dive,eq=all|eq=name|eq=description|eq=source"`
	PatternMatch string   `validate:"eq=exact|eq=partial"`
	Limit        int      `validate:"gte=1,lte=100"`
	Offset       int      `validate:"gte=0"`
}

type PostUserReq struct {
	AccountID    string `json:"account_id" db:"account_id" validate:"required,printascii"`
	EmailAddress string `json:"email_address" db:"email_address" validate:"required,email"`
	Password     string `json:"password" db:"password" validate:"required,printascii"`
	Name         string `json:"name" db:"name"`
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
