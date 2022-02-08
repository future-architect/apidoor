package model

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {

	fieldsSample := Fields{
		{
			Template: NewURITemplate("/normal"),
			Path:     NewURITemplate("/normal/dist"),
			Num:      5,
			Max:      10,
		},
		{
			Template: NewURITemplate("/exceeded"),
			Path:     NewURITemplate("/exceeded/dist"),
			Num:      5,
			Max:      2,
		},
		{
			Template: NewURITemplate("/unlimited"),
			Path:     NewURITemplate("/unlimited/dist"),
			Num:      5,
			Max:      "-",
		},
		{
			Template: NewURITemplate("/error/type"),
			Path:     NewURITemplate("/error/type/dist"),
			Num:      5,
			Max:      false,
		},
		{
			Template: NewURITemplate("/error/value"),
			Path:     NewURITemplate("/error/value/dist"),
			Num:      5,
			Max:      "unlimited",
		},
	}

	checkAPILimitTestData := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "valid request",
			input: "normal",
			err:   nil,
		},
		{
			name:  "limit exceeded",
			input: "exceeded",
			err:   errors.New("limit exceeded"),
		},
		{
			name:  "valid request",
			input: "unlimited",
			err:   nil,
		},
		{
			name:  "invalid type in limit",
			input: "error/type",
			err:   errors.New("unexpected limit value"),
		},
		{
			name:  "invalid value in limit",
			input: "error/value",
			err:   errors.New("unexpected limit value"),
		},
		{
			name:  "invalid path",
			input: "invalid",
			err:   nil,
		},
	}

	for _, tt := range checkAPILimitTestData {
		err := fieldsSample.CheckAPILimit(tt.input)

		if tt.err == nil && err != nil {
			t.Fatalf("case %s: unexpected error %s, want nil", tt.name, err.Error())
		}
		if tt.err != nil && err == nil {
			t.Fatalf("case %s: unexpected error nil, want %s", tt.name, tt.err.Error())
		}
		if tt.err != nil && err != nil {
			if tt.err.Error() != err.Error() {
				t.Fatalf("case %s: unexpected error %s, want %s", tt.name, err.Error(), tt.err.Error())
			}
		}
	}
}

func TestFields_URI(t *testing.T) {
	fields := Fields{
		{
			Template: NewURITemplate("/users/{userID}"),
			Path:     NewURITemplate("example.com/users/{userID}"),
		},
		{
			Template: NewURITemplate("/users"),
			Path:     NewURITemplate("example.com/users"),
		},
		{
			Template: NewURITemplate("/items/{itemID}"),
			Path:     NewURITemplate("example.com/v1/items/{itemID}"),
		},
		{
			Template: NewURITemplate("/items/{itemID}/{subItemID}"),
			Path:     NewURITemplate("example.com/items/{itemID}/{subItemID}"),
		},
		{
			Template: NewURITemplate("/{itemID}/old"),
			Path:     NewURITemplate("example.com/v1/{itemID}/old"),
		},
	}

	tests := []struct {
		name           string
		path           string
		wantForwardURI string
		wantErr        error
	}{
		{
			name:           "get path without placeholder properly",
			path:           "/users",
			wantForwardURI: "example.com/users",
			wantErr:        nil,
		},
		{
			name:           "get path with a placeholder properly",
			path:           "/users/foo",
			wantForwardURI: "example.com/users/foo",
			wantErr:        nil,
		},
		{
			name:           "get path whose placeholder's position is different properly",
			path:           "/items/bar",
			wantForwardURI: "example.com/v1/items/bar",
			wantErr:        nil,
		},
		{
			name:           "get path with multiple placeholders properly",
			path:           "/items/foo/bar",
			wantForwardURI: "example.com/items/foo/bar",
			wantErr:        nil,
		},
		{
			name:           "getting path when a placeholder is in the middle",
			path:           "/foo/old",
			wantForwardURI: "example.com/v1/foo/old",
			wantErr:        nil,
		},
		{
			name:           "get path fails when the path parameter is not provided",
			path:           "/items",
			wantForwardURI: "",
			wantErr:        ErrUnauthorizedRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fields.LookupTemplate(tt.path)
			if result == nil {
				if tt.wantForwardURI != "" {
					t.Error("returning non nil result is exprect, got nil")
				}
			} else if result.ForwardURL != tt.wantForwardURI {
				t.Errorf("uri in the result defers: want %s, got %s", tt.wantForwardURI, result.ForwardURL)
			}

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("returned error expects nil, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("returned error expects %v, got nil", tt.wantErr)
				} else if !errors.Is(err, tt.wantErr) {
					t.Errorf("returned error expects %v, got %v", tt.wantErr, err)

				}

			}

		})
	}
}
