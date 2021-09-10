package gateway_test

import (
	"encoding/csv"
	"gateway"
	"os"
	"testing"
)

func TestUpdateLog(t *testing.T) {
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}

	// TODO: modify test
	for i := 0; i < 2; i++ {
		gateway.UpdateLog("key", "path", nil)
	}

	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatal(err)
		}
		if line[1] != "key" {
			t.Fatalf("unexpected log %s, expected 'key'", line[1])
		} else if line[2] != "path" {
			t.Fatalf("unexpected log %s, expected 'path'", line[2])
		}
	}
}
