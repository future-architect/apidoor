package main

import (
	"log"
	"os"
	"strconv"
)

func GetEnvWithDeault(env string, def int) int {
	if def <= 0 {
		log.Fatal("default value should be positive")
	}

	n, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		log.Fatal(err)
	}
	if n <= 0 {
		n = def
	}

	return n
}
