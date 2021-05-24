package apidoor

import (
	"encoding/json"
	"log"
	"os"
)

var urldata OuterUrlData
var keydata OuterKeyData

func contains(list []int, a int) bool {
	for _, v := range list {
		if v == a {
			return true
		}
	}
	return false
}

func GetAPINum(url string) (int, error) {
	for i, v := range urldata.Url {
		if v == url {
			return i, nil
		}
	}

	return -1, &MyError{Message: "API not found"}
}

func RequestChecker(num int, key string) error {
	apilist, ok := keydata.Keys[key]
	if !ok {
		return &MyError{Message: "invalid key"}
	}

	if contains(apilist, num) {
		return nil
	}

	return &MyError{Message: "unauthorized request"}
}

func init() {
	urlfile, err := os.ReadFile("../../urlData.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(urlfile, &urldata); err != nil {
		log.Fatal(err)
	}

	keyfile, err := os.ReadFile("../../keyData.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(keyfile, &keydata); err != nil {
		log.Fatal(err)
	}
}
