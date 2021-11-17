package managementapi

import (
	"fmt"
	"gopkg.in/go-playground/validator.v8"
)

var (
	validate *validator.Validate
)

func init() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)
}

type BadRequestResp struct {
	Message          string           `json:"message"`
	InputValidations ValidationErrors `json:"input_validations,omitempty"`
}

type ValidationError struct {
	Field          string   `json:"field"`
	ConstraintType string   `json:"constraint_type"`
	Message        string   `json:"message"`
	Min            int      `json:"min,omitempty"`
	Max            int      `json:"max,omitempty"`
	Enum           []string `json:"enum,omitempty"`
	Got            string   `json:"got,omitempty"`
	Detail         string   `json:"detail,omitempty"`
}

func NewValidationError(fieldErr *validator.FieldError) *ValidationError {
	// TODO: impl
	ve := &ValidationError{}

	return ve
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("field: %s, constraint type: %s, message: %s", ve.Field, ve.ConstraintType, ve.Message)
}

func (ve ValidationError) generateMessage() string {
	msg := fmt.Sprintf("%s field biolates the following constraint: ", ve.Field)

	if ve.ConstraintType == "required" {
		msg += "required field, but got empty"
	}
	// TODO: FieldErrorからErrorの生成
	return msg

}

// ValidationErrorCustom allows for overriding output error message
type ValidationErrorCustom interface {
	customValidationError(fieldError *validator.FieldError) (*ValidationError, bool)
}

type ValidationErrors []*ValidationError

func (ves ValidationErrors) Error() string {
	ret := fmt.Sprintf("%d input validation(s) failed: [\n", len(ves))
	for _, v := range ves {
		ret += fmt.Sprintf("    %v,\n", v.Error())
	}
	ret += "]"

	return ret
}

func NewValidationErrors(target interface{}, fieldErrs validator.ValidationErrors) ValidationErrors {
	valErrs := ValidationErrors{}
	useCustomValidationError := false
	customErrorGenerator, ok := target.(ValidationErrorCustom)
	if ok {
		useCustomValidationError = true
	}

	for _, v := range fieldErrs {
		if useCustomValidationError {
			customErr, ok := customErrorGenerator.customValidationError(v)
			if ok {
				valErrs = append(valErrs, customErr)
				continue
			}
		}
		valErr := NewValidationError(v)
		valErrs = append(valErrs, valErr)
	}
	return valErrs
}

/*
type ValidationError struct {
	fieldErr *validator.FieldError
}

func NewValidationError(fieldErr *validator.FieldError) *ValidationError {
	return &ValidationError{
		fieldErr: fieldErr,
	}
}

func (ve ValidationError) Error() string {
	// TODO: フィールド名をjsonでのフィールド名に準拠 ex: TargetFieldではなく、target_field (SearchProductParams)
	msg := fmt.Sprintf("%s field biolates the following constraint: ", ve.fieldErr.NameNamespace)

	if ve.fieldErr.Tag == "required" {
		msg += "required field, but got empty"
	}
	// TODO: FieldErrorからErrorの生成
	return msg
}



type ValidationErrors struct {
	target interface{}
	errs []*ValidationError
}

func NewValidationErrors (target interface{}, fieldErrs validator.ValidationErrors) ValidationErrors {
	valErrs := ValidationErrors{
		target: target,
		errs: []*ValidationError{},
	}
	for _, v := range fieldErrs {
		valErr := NewValidationError(v)
		valErrs.errs = append(valErrs.errs, valErr)
	}
	return valErrs
}

func (ves ValidationErrors) Error() string {
	ret := fmt.Sprintf("%d input validation(s) failed: [\n", len(ves.errs))
	if vec, ok := ves.target.(ValidationErrorCustom); ok {
		for _, v := range ves.errs {
			ret += fmt.Sprintf("    %v,\n", vec.customError(v.fieldErr))
		}
	} else {
		for _, v := range ves.errs {
			ret += fmt.Sprintf("    %v,\n", v)
		}
	}
	ret += "]"

	return ret
}
*/

// ValidateStruct executes Validate.Struct method in go-playground/validator,
// and generates a validation error pointing out how input violates constraints.
// it returns nil or *ValidationErrors
func ValidateStruct(target interface{}) error {
	err := validate.Struct(target)
	if err == nil {
		return nil
	}
	fieldErrs := err.(validator.ValidationErrors)

	valErrors := NewValidationErrors(target, fieldErrs)

	return valErrors
}
