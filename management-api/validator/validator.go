package validator

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

//TODO: delete ../validator*.go

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

// BadRequestResp is used for a body of 4xx response
type BadRequestResp struct {
	Message          string            `json:"message"`
	ValidationErrors *ValidationErrors `json:"validation_errors,omitempty"`
}

func NewBadRequestResp(msg string) BadRequestResp {
	return BadRequestResp{
		Message: msg,
	}
}

func (br BadRequestResp) WriteResp(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp, err := json.Marshal(br)
	if err != nil {
		return fmt.Errorf("response to json error: %w", err)
	}
	if _, err = w.Write(resp); err != nil {
		return fmt.Errorf("write resoponse body error: %w", err)
	}
	return nil
}

type fieldError validator.FieldError

type ValidationError struct {
	Field          string      `json:"field"`
	ConstraintType string      `json:"constraint_type"`
	Message        string      `json:"message"`
	Lte            string      `json:"lte,omitempty"`
	Gte            string      `json:"gte,omitempty"`
	Ne             string      `json:"ne,omitempty"`
	Enum           []string    `json:"enum,omitempty"`
	Got            interface{} `json:"got,omitempty"`
}

func NewValidationError(fieldErr fieldError) *ValidationError {
	ve := &ValidationError{}

	ve.Field = trimFirstNameSpace(fieldErr.Namespace())

	ve.ConstraintType = generateConstraintType(fieldErr)

	ve.updateAttributes(fieldErr)

	ve.updateMessage()

	return ve
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("field: %s, constraint type: %s, message: %s", ve.Field, ve.ConstraintType, ve.Message)
}

func (ve *ValidationError) updateAttributes(fieldErr fieldError) {
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
	if ve.ConstraintType == "ne" || ve.ConstraintType == "length_ne" {
		ve.Ne = fieldErr.Param()
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
	case "ne":
		msg = fmt.Sprintf("input value is %v, but it must be not equal to %v", ve.Got, ve.Ne)
	case "length_gte":
		msg = fmt.Sprintf("input array length is %v, but it must be greater than or equal to %v", ve.Got, ve.Gte)
	case "length_lte":
		msg = fmt.Sprintf("input array length is %v, but it must be less than or equal to %v", ve.Got, ve.Lte)
	case "length_ne":
		msg = fmt.Sprintf("input array length is %v, but it must be not equal to %v", ve.Got, ve.Ne)
	default:
		msg = fmt.Sprintf("input value, %s, does not satisfy the format, %s", ve.Got, ve.ConstraintType)
	}

	ve.Message = msg
}

// trimFirstNameSpace trims the first element in nameSpace
// ex. User.address.city -> address.city
func trimFirstNameSpace(nameSpace string) string {
	nameSpaceList := strings.Split(nameSpace, ".")
	if len(nameSpaceList) < 2 {
		return nameSpace
	}

	return strings.Join(nameSpaceList[1:], ".")
}

func generateConstraintType(fieldErr validator.FieldError) string {
	// enum pattern, ex.) "eq=exact|eq=partial"
	if enumTagPattern.MatchString(fieldErr.Tag()) {
		return "enum"
	}

	if fieldErr.Kind() == reflect.Slice || fieldErr.Kind() == reflect.Map {
		if fieldErr.Tag() == "gte" {
			return "length_gte"
		} else if fieldErr.Tag() == "lte" {
			return "length_lte"
		} else if fieldErr.Tag() == "ne" {
			return "length_ne"
		}
	}

	return fieldErr.Tag()
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

func (ves *ValidationErrors) ToBadRequestResp() *BadRequestResp {
	resp := BadRequestResp{
		Message:          "input validation error",
		ValidationErrors: ves,
	}
	return &resp
}

func NewValidationErrors(fieldErrs validator.ValidationErrors) ValidationErrors {
	valErrs := ValidationErrors{}

	for _, v := range fieldErrs {
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

	valErrors := NewValidationErrors(fieldErrs)

	return valErrors
}
