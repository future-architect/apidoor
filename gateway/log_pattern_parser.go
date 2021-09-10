package gateway

import (
	"errors"
	"strings"
)

func LogPatternParser(pattern string) ([]string, error) {
	// parse input
	slice := strings.Split(pattern, "|")
	if len(slice) != 2 {
		return []string{}, errors.New("invalid use of separetor")
	}
	param, template := slice[0], slice[1]

	// name of header that user specify dynamically
	params := strings.Split(param, ",")
	// column name of log
	schema := strings.Split(template, ",")

	// allocate header name to schema
	usedParamNum := 0
	for i, v := range schema {
		if strings.HasPrefix(v, "%") {
			if len(params) <= usedParamNum {
				return []string{}, errors.New("the number of parameter does not match")
			}
			schema[i] = params[usedParamNum]
			usedParamNum++
		}
	}
	if usedParamNum != len(params) {
		return []string{}, errors.New("the number of parameter does not match")
	}

	return schema, nil
}
