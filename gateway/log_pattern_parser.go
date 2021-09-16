package gateway

import (
	"os"
	"strings"
)

var LogPattern []string

func init() {
	LogPattern = LogPatternParser()
}

func LogPatternParser() []string {
	// check log pattern
	env := os.Getenv("LOG_PATTERN")
	// column name of log
	pattern := strings.Split(env, ",")

	return pattern
}
