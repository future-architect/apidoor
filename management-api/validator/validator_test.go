package validator

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

type user struct {
	FirstName       string    `json:"first_name" validate:"required"`
	LastName        string    `json:"last_name" validate:"required"`
	BloodType       string    `json:"blood_type" validate:"required,eq=A|eq=B|eq=O|eq=AB"`
	Age             int       `json:"age" validate:"omitempty,gte=0,lte=130"`
	FavoriteNumbers []int     `json:"favorite_numbers" validate:"gte=1,dive,gte=1,lte=100"`
	Email           string    `json:"email" validate:"omitempty,email"`
	Addresses       []address `json:"addresses" validate:"omitempty,dive,required"`
}

type address struct {
	Country string `json:"country" validate:"required"`
	City    string `json:"city" validate:"required"`
}

func TestValidator(t *testing.T) {
	tests := []struct {
		name       string
		input      user
		wantErr    *ValidationErrors
		wantErrMsg string
	}{
		{
			name: "required field is empty",
			input: user{
				FirstName:       "foo",
				BloodType:       "AB",
				FavoriteNumbers: []int{12},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "last_name",
					ConstraintType: "required",
					Got:            "",
					Message:        `required field, but got empty`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: last_name, constraint type: required, message: required field, but got empty,
]`,
		},
		{
			name: "value violates enum constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "C",
				FavoriteNumbers: []int{12},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "blood_type",
					ConstraintType: "enum",
					Enum:           []string{"A", "B", "O", "AB"},
					Got:            "C",
					Message:        `input value is C, but it must be one of the following values: [A B O AB]`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: blood_type, constraint type: enum, message: input value is C, but it must be one of the following values: [A B O AB],
]`,
		},
		{
			name: "value violates lte constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				Age:             200,
				FavoriteNumbers: []int{12},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "age",
					ConstraintType: "lte",
					Lte:            "130",
					Got:            200,
					Message:        `input value is 200, but it must be less than or equal to 130`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: age, constraint type: lte, message: input value is 200, but it must be less than or equal to 130,
]`,
		},
		{
			name: "value violates gte_length constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "favorite_numbers",
					ConstraintType: "length_gte",
					Gte:            "1",
					Got:            0,
					Message:        `input array length is 0, but it must be greater than or equal to 1`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: favorite_numbers, constraint type: length_gte, message: input array length is 0, but it must be greater than or equal to 1,
]`,
		},
		{
			name: "value in slice violates gte constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{2, 0},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "favorite_numbers[1]",
					ConstraintType: "gte",
					Gte:            "1",
					Got:            0,
					Message:        `input value is 0, but it must be greater than or equal to 1`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: favorite_numbers[1], constraint type: gte, message: input value is 0, but it must be greater than or equal to 1,
]`,
		},
		{
			name: "email field's value does not satisfy email format",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				Email:           "foo.@example.com",
				FavoriteNumbers: []int{2, 1},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "email",
					ConstraintType: "email",
					Got:            "foo.@example.com",
					Message:        `input value, foo.@example.com, does not satisfy the format, email`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: email, constraint type: email, message: input value, foo.@example.com, does not satisfy the format, email,
]`,
		},
		{
			name: "email field's value does not satisfy email format",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{2, 1},
				Email:           "foo.@example.com",
			},
			wantErr: &ValidationErrors{
				{
					Field:          "email",
					ConstraintType: "email",
					Got:            "foo.@example.com",
					Message:        `input value, foo.@example.com, does not satisfy the format, email`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: email, constraint type: email, message: input value, foo.@example.com, does not satisfy the format, email,
]`,
		},
		{
			name: "validation of sub-struct fails",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{2, 1},
				Addresses: []address{
					{
						Country: "Japan",
					},
				},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "addresses[0].city",
					ConstraintType: "required",
					Got:            "",
					Message:        `required field, but got empty`,
				},
			},
			wantErrMsg: `1 input validation(s) failed: [
    field: addresses[0].city, constraint type: required, message: required field, but got empty,
]`,
		},
		{
			name: "multiple errors occur",
			input: user{
				FirstName:       "foo",
				BloodType:       "A",
				FavoriteNumbers: []int{2, 1, -2},
				Addresses: []address{
					{
						Country: "Japan",
					},
				},
			},
			wantErr: &ValidationErrors{
				{
					Field:          "last_name",
					ConstraintType: "required",
					Got:            "",
					Message:        `required field, but got empty`,
				},
				{
					Field:          "favorite_numbers[2]",
					ConstraintType: "gte",
					Gte:            "1",
					Got:            -2,
					Message:        `input value is -2, but it must be greater than or equal to 1`,
				},
				{
					Field:          "addresses[0].city",
					ConstraintType: "required",
					Got:            "",
					Message:        `required field, but got empty`,
				},
			},
			wantErrMsg: `3 input validation(s) failed: [
    field: last_name, constraint type: required, message: required field, but got empty,
    field: favorite_numbers[2], constraint type: gte, message: input value is -2, but it must be greater than or equal to 1,
    field: addresses[0].city, constraint type: required, message: required field, but got empty,
]`,
		},
		{
			name: "no validation error occurs",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{2, 1},
				Addresses: []address{
					{
						Country: "Japan",
						City:    "Tokyo",
					},
				},
			},
			wantErr:    nil,
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if err == nil {
				if tt.wantErr != nil {
					t.Errorf("no error returned, but it is expected that the following error occurs\n%v", tt.wantErr)
				}
				return
			}

			if tt.wantErr == nil {
				t.Errorf("no error occurs is expected, but got the following error\n%v", err)
				return
			}

			if diff := cmp.Diff(*tt.wantErr, err); diff != "" {
				t.Errorf("error differs,\n%v", diff)
			}

			if tt.wantErrMsg != err.Error() {
				t.Errorf("error message differs,\nwant: %v,\ngot: %v", tt.wantErrMsg, err.Error())
			}
		})
	}
}
