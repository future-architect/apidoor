package apidoor

import (
	"encoding/json"
	"log"
	"os"
)

var Urldata OuterUrlData
var Keydata OuterKeyData

func contains(list []int, a int) bool {
	for _, v := range list {
		if v == a {
			return true
		}
	}
	return false
}

func GetAPIURL(num int, key string) (string, error) {
	apilist, ok := Keydata.Keys[key]
	if !ok {
		return "", &MyError{Message: "invalid key"}
	}

	if !contains(apilist, num) {
		return "", &MyError{Message: "unauthorized request"}
	}

	return Urldata.Url[num], nil
}

func init() {
	urlfile, err := os.ReadFile(os.Getenv("GOPATH") + "/src/apidoor/urlData.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(urlfile, &Urldata); err != nil {
		log.Fatal(err)
	}

	keyfile, err := os.ReadFile(os.Getenv("GOPATH") + "/src/apidoor/keyData.json")
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(keyfile, &Keydata); err != nil {
		log.Fatal(err)
	}
}
