package model

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {

	cases := []struct {
		name   string
		input  string
		output string
		err    error
	}{
		{
			name:   "valid path",
			input:  "/normal",
			output: "normal/dist",
			err:    nil,
		},
		{
			name:   "invalid path",
			input:  "/invalid",
			output: "",
			err:    &MyError{Message: "unauthorized request"},
		},
	}

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

	for _, tt := range cases {
		path, err := fieldsSample.URI(tt.input)

		if tt.output != path {
			t.Fatalf("case %s: unexpected result %s, want %s", tt.name, path, tt.output)
		}
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
