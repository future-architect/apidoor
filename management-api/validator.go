package managementapi

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

var (
	validate *validator.Validate

	enumTagPattern    = regexp.MustCompile(`^(eq=\w+\|)*eq=\w+$`)
	enumTagSubPattern = regexp.MustCompile(`eq=(\w+)`)
)

func init() {
	validate = validator.New()

	// use "json" tag value
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonTag, ok := fld.Tag.Lookup("json")
		if !ok {
			// if "json" tag is not set, use struct field name.
			return fld.Name
		}

		// exclude second and later field, such as "omitempty"
		name := strings.SplitN(jsonTag, ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}

type BadRequestResp struct {
	Message          string            `json:"message"`
	InputValidations *ValidationErrors `json:"input_validations,omitempty"`
}

type ValidationError struct {
	Field          string      `json:"field"`
	ConstraintType string      `json:"constraint_type"`
	Message        string      `json:"message"`
	Lte            string      `json:"lte,omitempty"`
	Gte            string      `json:"gte,omitempty"`
	Enum           []string    `json:"enum,omitempty"`
	Got            interface{} `json:"got,omitempty"`
}

func NewValidationError(fieldErr validator.FieldError) *ValidationError {
	ve := &ValidationError{}

	ve.Field = fieldErr.Field()

	ve.ConstraintType = generateConstraintType(fieldErr)

	ve.updateAttributes(fieldErr)

	ve.updateMessage()

	return ve
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("field: %s, constraint type: %s, message: %s", ve.Field, ve.ConstraintType, ve.Message)
}

func (ve *ValidationError) updateAttributes(fieldErr validator.FieldError) {
	if ve.ConstraintType == "enum" {
		matches := enumTagSubPattern.FindAllStringSubmatch(fieldErr.Tag(), -1)
		ve.Enum = make([]string, len(matches))
		for i, v := range matches {
			ve.Enum[i] = v[1]
		}
	}
	if ve.ConstraintType == "lte" || ve.ConstraintType == "length_lte" {
		ve.Lte = fieldErr.Param()
	}
	if ve.ConstraintType == "gte" || ve.ConstraintType == "length_gte" {
		ve.Gte = fieldErr.Param()
	}

	if fieldErr.Kind() == reflect.Slice || fieldErr.Kind() == reflect.Map {
		ve.Got = reflect.ValueOf(fieldErr.Value()).Len()
	} else {
		ve.Got = fieldErr.Value()
	}

}

func (ve *ValidationError) updateMessage() {
	var msg string

	switch ve.ConstraintType {
	case "required":
		msg = "required field, but got empty"
	case "enum":
		msg = fmt.Sprintf("input value is %v, but it must be one of the following values: %v", ve.Got, ve.Enum)
	case "gte":
		msg = fmt.Sprintf("input value is %v, but it must be greater than or equal to %v", ve.Got, ve.Gte)
	case "lte":
		msg = fmt.Sprintf("input value is %v, but it must be less than or equal to %v", ve.Got, ve.Lte)
	case "length_gte":
		msg = fmt.Sprintf("input array length is %v, but it must be greater than or equal to %v", ve.Got, ve.Gte)
	case "length_lte":
		msg = fmt.Sprintf("input array length is %v, but it must be less than or equal to %v", ve.Got, ve.Lte)
	default:
		msg = fmt.Sprintf("input value, %s, does not satisfy the format, %s", ve.Got, ve.ConstraintType)
	}

	ve.Message = msg
}

func generateConstraintType(fieldErr validator.FieldError) string {
	// enum pattern, ex.) "eq=exact|eq=partial"
	if enumTagPattern.MatchString(fieldErr.Tag()) {
		return "enum"
	}

	fmt.Println(fieldErr.Kind())
	if fieldErr.Kind() == reflect.Slice || fieldErr.Kind() == reflect.Map {
		if fieldErr.Tag() == "gte" {
			return "length_gte"
		} else if fieldErr.Tag() == "lte" {
			return "length_lte"
		}
	}

	return fieldErr.Tag()
}

// ValidationErrorCustom allows for overriding output error message
type ValidationErrorCustom interface {
	customValidationError(fieldError validator.FieldError) (*ValidationError, bool)
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
