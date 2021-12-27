package model

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {
	type fieldsURITest struct {
		input  string
		output string
		err    error
	}

	type fieldsCheckAPILimitTest struct {
		input string
		err   error
	}

	cases := []fieldsURITest{
		// valid path
		{
			input:  "/normal",
			output: "normal/dist",
			err:    nil,
		},
		// invalid path
		{
			input:  "/invalid",
			output: "",
			err:    &MyError{Message: "unauthorized request"},
		},
	}

	var checkAPILimitTestData = []fieldsCheckAPILimitTest{
		// valid request
		{
			input: "normal",
			err:   nil,
		},
		// limit exceeded
		{
			input: "exceeded",
			err:   errors.New("limit exceeded"),
		},
		// valid request
		{
			input: "unlimited",
			err:   nil,
		},
		// invalid type in limit
		{
			input: "error/type",
			err:   errors.New("unexpected limit value"),
		},
		// invalid value in limit
		{
			input: "error/value",
			err:   errors.New("unexpected limit value"),
		},
		// invalid path
		{
			input: "invalid",
			err:   nil,
		},
	}

	var fieldsSample = Fields{
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

	for i, tt := range cases {
		path, err := fieldsSample.URI(tt.input)

		if tt.output != path {
			t.Fatalf("case %d: unexpected result %s, want %s", i, path, tt.output)
		}
		if tt.err == nil && err != nil {
			t.Fatalf("case %d: unexpected error %s, want nil", i, err.Error())
		}
		if tt.err != nil && err == nil {
			t.Fatalf("case %d: unexpected error nil, want %s", i, tt.err.Error())
		}
		if tt.err != nil && err != nil {
			if tt.err.Error() != err.Error() {
				t.Fatalf("case %d: unexpected error %s, want %s", i, err.Error(), tt.err.Error())
			}
		}
	}

	for i, tt := range checkAPILimitTestData {
		err := fieldsSample.CheckAPILimit(tt.input)

		if tt.err == nil && err != nil {
			t.Fatalf("case %d: unexpected error %s, want nil", i, err.Error())
		}
		if tt.err != nil && err == nil {
			t.Fatalf("case %d: unexpected error nil, want %s", i, tt.err.Error())
		}
		if tt.err != nil && err != nil {
			if tt.err.Error() != err.Error() {
				t.Fatalf("case %d: unexpected error %s, want %s", i, err.Error(), tt.err.Error())
			}
		}
	}
}
