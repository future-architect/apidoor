package gateway

import (
	"os"
	"strings"
)

var LogPattern []LogElement

func init() {
	Parse()
}

func Parse() {
	// check log pattern
	env := os.Getenv("LOG_PATTERN")

	// get column name of log
	var pattern []LogElement
	for _, value := range strings.Split(env, ",") {
		switch value {
		case "time":
			pattern = append(pattern, TimeElement{})
		case "key":
			pattern = append(pattern, APIKeyElement{})
		case "path":
			pattern = append(pattern, PathElement{})
		default:
			pattern = append(pattern, HeaderElement{value})
		}
	}
	LogPattern = pattern
}
