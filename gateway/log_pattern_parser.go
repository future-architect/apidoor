package gateway

import (
	"os"
	"strings"
)

var LogPattern []LogElement

func init() {
	LogPattern = LogPatternParser()
}

func LogPatternParser() []LogElement {
	// check log pattern
	env := os.Getenv("LOG_PATTERN")

	// get column name of log
	var pattern []LogElement
	for _, value := range strings.Split(env, ",") {
		pattern = append(pattern, LogElement(value))
	}

	return pattern
}
