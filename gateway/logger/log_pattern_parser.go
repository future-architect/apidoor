package logger

import (
	"os"
	"strings"
)

var LogOptionPattern []LogOption

func init() {
	LogOptionPattern = LogPatternParser()
}

func LogPatternParser() []LogOption {
	// check log pattern
	env := os.Getenv("LOG_PATTERN")

	// get column name and set function to write log
	var pattern []LogOption
	for _, value := range strings.Split(env, ",") {
		switch value {
		case "time":
			pattern = append(pattern, WithTime())
		case "key":
			pattern = append(pattern, WithKey())
		case "path":
			pattern = append(pattern, WithPath())
		default:
			pattern = append(pattern, HeaderElement(value))
		}
	}

	return pattern
}
