package swaggerparser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"net/url"
)

type SwaggerVersion string

var (
	SwaggerV2 SwaggerVersion = "2"
	SwaggerV3 SwaggerVersion = "3"
)

const (
	basePathFieldName = "x-apidoor-base-path"
	pathFieldName     = "x-apidoor-path"
)

var supportedSchemes = []string{"https", "http"}

type Swagger struct {
	Version        SwaggerVersion
	Schemes        []string
	ForwardURLBase string
	PathBase       string
	APIs           []API
}

type API struct {
	ForwardURL string
	Path       string
}

type Parser struct {
	Fetcher FileFetcher

	data map[string]interface{}
	url  *url.URL
}

func NewParser(fetcher FileFetcher) Parser {
	return Parser{
		Fetcher: fetcher,
	}
}

/*
Parse fetches a swagger file over the network and parse the file as a swagger definition file based on the following rules
swagger v2: https://swagger.io/docs/specification/2-0/basic-structure/
swagger v3: https://swagger.io/docs/specification/basic-structure/
*/
func (sp *Parser) Parse(ctx context.Context, swaggerUrl string) (*Swagger, error) {
	var err error
	if sp.url, err = url.Parse(swaggerUrl); err != nil {
		return nil, newError(FormatError, fmt.Errorf("requested swagger url is not url format, %s", swaggerUrl))
	}

	file, err := sp.Fetcher.Fetch(ctx, sp.url)
	if err != nil {
		return nil, err
	}
	return sp.parse(file)
}

func (sp *Parser) parse(swag *swaggerFile) (*Swagger, error) {
	sp.data = make(map[string]interface{})

	switch swag.contentType {
	case typeJson:
		if err := json.Unmarshal(swag.body.Bytes(), &sp.data); err != nil {
			return nil, newError(FileParseError, errors.New("cannot parse data as json"))
		}
	case typeYaml:
		if err := yaml.Unmarshal(swag.body.Bytes(), &sp.data); err != nil {
			return nil, newError(FileParseError, errors.New("cannot parse data as yaml"))
		}
	}

	version := sp.getSwaggerVersion()
	if version == nil {
		return nil, newError(FileParseError, errors.New("the swagger file is not based on v2 nor v3"))
	}

	switch *version {
	case SwaggerV2:
		return newParserV2(sp).parse()
	case SwaggerV3:
		return newParserV3(sp).parse()
	}
	// unreachable
	return nil, newError(OtherError, errors.New("parser is not implemented"))
}

func (sp Parser) getSwaggerVersion() *SwaggerVersion {
	if v, ok := sp.data["swagger"]; ok {
		if v, ok := v.(string); ok {
			if v == "2.0" {
				return &SwaggerV2
			}
			return nil
		}
		return nil
	}

	if v, ok := sp.data["openapi"]; ok {
		if v, ok := v.(string); ok {
			switch v {
			case "3.0.0", "3.0.1", "3.0.2", "3.0.3":
				return &SwaggerV3
			default:
				return nil
			}
		}
		return nil
	}
	return nil
}
