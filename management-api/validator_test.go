package managementapi_test

import (
	"github.com/future-architect/apidoor/managementapi"
	"testing"
)

// User contains user information
type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
}

// TODO: アップデート
func TestValidator(t *testing.T) {
	tests := []struct {
		name    string
		input   User
		wantErr string
	}{
		{
			name: "required fieldが空",
			input: User{
				FirstName: "hoge",
			},
			wantErr: `1 input validation(s) failed: [
    LastName field biolates the following constraint: required field, but got empty,
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := managementapi.ValidateStruct(tt.input)
			if err == nil {
				if tt.wantErr != "" {
					t.Errorf("expected nil error, got %s", tt.wantErr)
				}
				return
			}

			if err.Error() != tt.wantErr {
				t.Errorf("error msg differs: want %s, got %s", err.Error(), tt.wantErr)
			}
		})
	}

}
