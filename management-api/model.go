package managementapi

import (
	"gopkg.in/go-playground/validator.v8"

	"github.com/gorilla/schema"
)

var (
	validate      *validator.Validate
	schemaDecoder *schema.Decoder
)

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)

	schemaDecoder = schema.NewDecoder()
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
	Q            string `schema:"name" validate:"required"`
	TargetFields string `schema:"target_fields"`
	PatternMatch string `schema:"pattern_match"`
	Limit        int    `schema:"limit" validate:"gte=1,lte=100"`
	Offset       int    `schema:"offset"`
}

type SearchProductsResp struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
}
