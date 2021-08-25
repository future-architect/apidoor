package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/welcome", handler)
	http.ListenAndServe(":3333", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "welcome to apidoor!\n")
}
