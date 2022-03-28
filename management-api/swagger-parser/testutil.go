package swaggerparser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

type TestFetcher struct{}

func (tf TestFetcher) Fetch(ctx context.Context, url *url.URL) (*swaggerFile, error) {
	filePath := ""
	switch url.String() {
	// swagger 2.0
	case "http://api.example.com/v2/swagger.json":
		filePath = "./testdata/testv2.json"
	case "http://api.example.com/v2/no_host_provided/swagger.json":
		filePath = "./testdata/testv2_no_host_provided.json"
	case "http://api.example.com/v2/wrong_format/swagger.json":
		filePath = "./testdata/testv2_wrong_format.json"
	// openapi 3.0
	case "http://api.example.com/v3/swagger.yaml":
		filePath = "./testdata/testv3.yaml"
	case "http://api.example.com/v3/no_servers_provided/swagger.yaml":
		filePath = "./testdata/testv3_no_servers_provided.yaml"
	case "http://api.example.com/v3/wrong_format/swagger.yaml":
		filePath = "./testdata/testv3_wrong_format.yaml"
	case "http://api.example.com/v4/swagger.yaml":
		filePath = "./testdata/testv4.yaml"
	default:
		return nil, newError(FormatError, NotFoundErr)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, newErrorString(OtherError, "no such definition file")
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, file); err != nil {
		return nil, newError(OtherError, fmt.Errorf("io copy failed: %w", err))
	}

	var contentType contentType
	if strings.HasSuffix(filePath, "json") {
		contentType = typeJson
	} else if strings.HasSuffix(filePath, "yaml") {
		contentType = typeYaml
	}

	return &swaggerFile{
		body:        buf,
		contentType: contentType,
	}, nil
}
