package swagger_parser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"
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

func TestParser(t *testing.T) {
	parser := NewParser(TestFetcher{})

	tests := []struct {
		name        string
		urlStr      string
		wantSwagger *Swagger
		wantErr     *Error
	}{
		{
			name:   "parse swagger v2 json file",
			urlStr: "http://api.example.com/v2/swagger.json",
			wantSwagger: &Swagger{
				Version:        "v2",
				Schemes:        []string{"https"},
				ForwardURLBase: "api.example.com/sample",
				PathBase:       "/sample_gateway",
				APIs: []API{
					{
						ForwardURL: "/users",
						Path:       "/sample_users",
					},
					{
						ForwardURL: "/users/{user_id}",
						Path:       "/users/{user_id}",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "parse swagger v2 json file that omits host field",
			urlStr: "http://api.example.com/v2/no_host_provided/swagger.json",
			wantSwagger: &Swagger{
				Version:        "v2",
				Schemes:        []string{"https"},
				ForwardURLBase: "api.example.com/sample",
				PathBase:       "/sample_gateway",
				APIs: []API{
					{
						ForwardURL: "/users",
						Path:       "/sample_users",
					},
					{
						ForwardURL: "/users/{user_id}",
						Path:       "/users/{user_id}",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:        "error occurs when target v2 json file format is wrong",
			urlStr:      "http://api.example.com/v2/wrong_format/swagger.json",
			wantSwagger: nil,
			wantErr: &Error{
				ErrorType: FileParseError,
				Message:   errors.New("the value of the host field is not string"),
			},
		},
		{
			name:   "parse swagger v3 yaml file",
			urlStr: "http://api.example.com/v3/swagger.yaml",
			wantSwagger: &Swagger{
				Version:        "v3",
				Schemes:        []string{"https"},
				ForwardURLBase: "api.example.com/v3",
				PathBase:       "/base",
				APIs: []API{
					{
						ForwardURL: "/users",
						Path:       "/foo",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "parse swagger v3 yaml file that omits servers field",
			urlStr: "http://api.example.com/v3/no_servers_provided/swagger.yaml",
			wantSwagger: &Swagger{
				Version:        "v3",
				Schemes:        []string{"http"},
				ForwardURLBase: "api.example.com/v3/no_servers_provided",
				PathBase:       "/base",
				APIs: []API{
					{
						ForwardURL: "/users",
						Path:       "/users",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:        "error occurs when target v3 yaml file format is wrong",
			urlStr:      "http://api.example.com/v3/wrong_format/swagger.yaml",
			wantSwagger: nil,
			wantErr: &Error{
				ErrorType: FileParseError,
				Message:   errors.New("x-apidoor-base-path field must be string field"),
			},
		},
		{
			name:        "unsupported version file is provided",
			urlStr:      "http://api.example.com/v4/swagger.yaml",
			wantSwagger: nil,
			wantErr: &Error{
				ErrorType: FileParseError,
				Message:   errors.New("the swagger file is not based on v2 nor v3"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			swagger, err := parser.Parse(context.Background(), tt.urlStr)

			if diff := cmp.Diff(tt.wantSwagger, swagger); diff != "" {
				t.Errorf("returned swagger differs:\n%s", diff)
			}

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("no error returned is expected, got %v", err)
				}
				return
			}
			if diff := cmp.Diff(tt.wantErr.Error(), err.Error()); diff != "" {
				t.Errorf("returned error differs:\n%s", diff)
			}

		})
	}

}
