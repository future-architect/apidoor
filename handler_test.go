package apidoor_test

import (
	"apidoor"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var table = []int{
	http.StatusOK,
	http.StatusBadRequest,
	http.StatusInternalServerError,
}

func TestHandler(t *testing.T) {
	for _, code := range table {
		message := []byte("response from API server")
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			w.Write(message)
		}))
		defer ts.Close()

		// change destination of API temporarily
		apidoor.Urldata.Url = []string{
			ts.URL,
		}

		r := httptest.NewRequest(http.MethodGet, "/hoge?apikey=apikey1&num=0", nil)
		w := httptest.NewRecorder()
		apidoor.Handler(w, r)

		rw := w.Result()
		defer rw.Body.Close()

		if rw.StatusCode != code {
			t.Fatalf("unexpected status code %d, expected %d", rw.StatusCode, code)
		}

		b, err := io.ReadAll(rw.Body)
		if err != nil {
			t.Fatal("unexpected body type")
		}

		trimmed := strings.TrimSpace(string(b))
		if trimmed != string(message) {
			t.Fatalf("unexpected response: %s, expected: %s", trimmed, string(message))
		}

	}
}
