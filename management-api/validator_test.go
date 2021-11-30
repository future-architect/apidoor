package managementapi_test

import (
	"github.com/future-architect/apidoor/managementapi"
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
		name    string
		input   user
		wantErr *managementapi.ValidationErrors
	}{
		{
			name: "required field is empty",
			input: user{
				FirstName:       "foo",
				BloodType:       "AB",
				FavoriteNumbers: []int{12},
			},
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "last_name",
					ConstraintType: "required",
					Got:            "",
					Message:        `required field, but got empty`,
				},
			},
		},
		{
			name: "value violates enum constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "C",
				FavoriteNumbers: []int{12},
			},
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "blood_type",
					ConstraintType: "enum",
					Enum:           []string{"A", "B", "O", "AB"},
					Got:            "C",
					Message:        `input value is C, but it must be one of the following values: [A B O AB]`,
				},
			},
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
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "age",
					ConstraintType: "lte",
					Lte:            "130",
					Got:            200,
					Message:        `input value is 200, but it must be less than or equal to 130`,
				},
			},
		},
		{
			name: "value violates gte_length constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{},
			},
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "favorite_numbers",
					ConstraintType: "length_gte",
					Gte:            "1",
					Got:            0,
					Message:        `input array length is 0, but it must be greater than or equal to 1`,
				},
			},
		},
		{
			name: "value in slice violates gte constraint",
			input: user{
				FirstName:       "foo",
				LastName:        "bar",
				BloodType:       "A",
				FavoriteNumbers: []int{2, 0},
			},
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "favorite_numbers[1]",
					ConstraintType: "gte",
					Gte:            "1",
					Got:            0,
					Message:        `input value is 0, but it must be greater than or equal to 1`,
				},
			},
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
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "email",
					ConstraintType: "email",
					Got:            "foo.@example.com",
					Message:        `input value, foo.@example.com, does not satisfy the format, email`,
				},
			},
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
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "email",
					ConstraintType: "email",
					Got:            "foo.@example.com",
					Message:        `input value, foo.@example.com, does not satisfy the format, email`,
				},
			},
		},
		{
			//
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
			wantErr: &managementapi.ValidationErrors{
				{
					Field:          "city",
					ConstraintType: "required",
					Got:            "",
					Message:        `required field, but got empty`,
				},
			},
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
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := managementapi.ValidateStruct(tt.input)
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
		})
	}
}
