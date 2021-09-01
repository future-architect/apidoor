package gateway_test

import (
	"errors"
	"gateway"
	"testing"
)

type fieldsURITest struct {
	input  string
	output string
	err    error
}

type fieldsCheckAPILimitTest struct {
	input string
	err   error
}

var URITestData = []fieldsURITest{
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
		err:    &gateway.MyError{Message: "unauthorized request"},
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

var fieldsSample = gateway.Fields{
	{
		Template: *gateway.NewURITemplate("/normal"),
		Path:     *gateway.NewURITemplate("/normal/dist"),
		Num:      5,
		Max:      10,
	},
	{
		Template: *gateway.NewURITemplate("/exceeded"),
		Path:     *gateway.NewURITemplate("/exceeded/dist"),
		Num:      5,
		Max:      2,
	},
	{
		Template: *gateway.NewURITemplate("/unlimited"),
		Path:     *gateway.NewURITemplate("/unlimited/dist"),
		Num:      5,
		Max:      "-",
	},
	{
		Template: *gateway.NewURITemplate("/error/type"),
		Path:     *gateway.NewURITemplate("/error/type/dist"),
		Num:      5,
		Max:      false,
	},
	{
		Template: *gateway.NewURITemplate("/error/value"),
		Path:     *gateway.NewURITemplate("/error/value/dist"),
		Num:      5,
		Max:      "unlimited",
	},
}

func TestFields(t *testing.T) {
	for i, tt := range URITestData {
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
