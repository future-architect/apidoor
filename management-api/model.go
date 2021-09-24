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
	Name          string `schema:"name" db:"name"`
	IsNameExact   bool   `schema:"is_name_exact" db:"is_name_exact"`
	Source        string `schema:"source" db:"source"`
	IsSourceExact bool   `schema:"is_source_exact" db:"is_source_exact"`
	Description   string `schema:"description" db:"description"`
	Keyword       string `schema:"keyword" db:"keyword"`
}
