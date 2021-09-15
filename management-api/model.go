package managementapi

import "errors"

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
	Name        string `json:"name" db:"name"`
	Source      string `json:"source" db:"source"`
	Description string `json:"description" db:"description"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail"`
	SwaggerURL  string `json:"swagger_url" db:"swagger_url"`
}

func (pp PostProductReq) CheckNoEmptyField() error {
	if pp.Name == "" {
		return errors.New("name field required")
	}
	if pp.Source == "" {
		return errors.New("source field required")
	}
	if pp.Description == "" {
		return errors.New("description field required")
	}
	if pp.Thumbnail == "" {
		return errors.New("thumbnail field required")
	}
	if pp.SwaggerURL == "" {
		return errors.New("swagger_url field required")
	}
	return nil
}
