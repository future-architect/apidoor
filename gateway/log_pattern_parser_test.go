package gateway_test

import (
	"errors"
	"gateway"
	"testing"
)

type logPatternParserTest struct {
	input  string
	output []string
	err    error
}

var logPatternParserTestData = []logPatternParserTest{
	{
		input: "TEST1,TEST2|log,%x,%x",
		output: []string{
			"log",
			"TEST1",
			"TEST2",
		},
		err: nil,
	},
	{
		input:  "TEST1,TEST2,TEST3|log,%x,%x",
		output: []string{},
		err:    errors.New("the number of parameter does not match"),
	},
	{
		input:  "TEST1|log,%x,%x",
		output: []string{},
		err:    errors.New("the number of parameter does not match"),
	},
	{
		input:  "TEST1|TEST1|log,%x,%x",
		output: []string{},
		err:    errors.New("invalid use of separetor"),
	},
}

func TestLogPatternParser(t *testing.T) {
	for index, tt := range logPatternParserTestData {
		schema, err := gateway.LogPatternParser(tt.input)

		// test if the length of output is valid
		if len(tt.output) != len(schema) {
			t.Fatalf("case %d: unexpected length of output %d, expected %d", index, len(schema), len(tt.output))
		}

		// test if the output value is valid
		for i, v := range tt.output {
			if schema[i] != v {
				t.Fatalf("case %d: unexpected output %s, expected %s", index, schema[i], v)
			}
		}

		// test if the error message is valid
		if tt.err == nil && err != nil {
			t.Fatalf("case %d: unexpected error %s, want nil", index, err.Error())
		}
		if tt.err != nil && err == nil {
			t.Fatalf("case %d: unexpected error nil, want %s", index, tt.err.Error())
		}
		if tt.err != nil && err != nil {
			if tt.err.Error() != err.Error() {
				t.Fatalf("case %d: unexpected error %s, want %s", index, err.Error(), tt.err.Error())
			}
		}
	}
}
