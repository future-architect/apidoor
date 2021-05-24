package apidoor

import "fmt"

type MyError struct {
	Message string `json:"message"`
}

func (err *MyError) Error() string {
	return fmt.Sprintf("error: %s", err.Message)
}

type OuterUrlData struct {
	Url []string `json:"url"`
}

type OuterKeyData struct {
	Keys map[string][]int `json:"keys"`
}
