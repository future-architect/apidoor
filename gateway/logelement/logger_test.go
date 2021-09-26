package gateway_test

import (
	gateway "gateway/logelement"
	"net/http"
	"os"
	"testing"
)

func TestUpdateLog(t *testing.T) {

	// save current environment variable
	tmp := os.Getenv("LOG_PATTERN")
	t.Cleanup(func() {
		os.Setenv("LOG_PATTERN", tmp)
	})

	os.Setenv("LOG_PATTERN", "time,key,path,TEST1,TEST2")
	gateway.Parse()

	// run test
	r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Set("TEST1", "header1")
	r.Header.Set("TEST2", "header2")

	for i := 0; i < 2; i++ {
		gateway.UpdateLog("key", "path", r)
	}
}
