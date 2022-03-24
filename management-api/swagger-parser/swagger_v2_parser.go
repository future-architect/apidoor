package swagger_parser

import (
	"fmt"
	"strings"
)

type parserV2 struct {
	*Parser
}

func newParserV2(base *Parser) *parserV2 {
	return &parserV2{
		base,
	}
}

func (p *parserV2) parse() (*Swagger, error) {
	forwardURLBase, err := p.getForwardBase()
	if err != nil {
		return nil, err
	}

	schemes, err := p.getSchemes()
	if err != nil {
		return nil, err
	}

	pathBase, err := p.getBaseGatewayPath()
	if err != nil {
		return nil, err
	}

	apis, err := p.parsePaths()
	if err != nil {
		return nil, err
	}

	return &Swagger{
		Version:        "v2",
		Schemes:        schemes,
		ForwardURLBase: forwardURLBase,
		PathBase:       pathBase,
		APIs:           apis,
	}, nil
}

func (p parserV2) getForwardBase() (string, error) {
	host, err := p.getHost()
	if err != nil {
		return "", err
	}
	basePath, err := p.getBaseForwardPath()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", host, basePath), nil
}

func (p parserV2) getHost() (string, error) {
	hostField, ok := p.data["host"]
	if ok {
		if host, ok := hostField.(string); ok {
			if isOnlyHostContained(host) {
				return host, nil
			}
			return "", newErrorString(FileParseError, "host field must be contain hostname and port (optional)")
		}
		return "", newErrorString(FileParseError, "the value of the host field is not string")
	}
	return p.url.Host, nil
}

func (p parserV2) getBaseForwardPath() (string, error) {
	basePathField, ok := p.data["basePath"]
	if ok {
		if basePath, ok := basePathField.(string); ok {
			if strings.HasPrefix(basePath, "/") {
				return basePath, nil
			}
			return "", newErrorString(FileParseError, "base path field must start with '/'")
		}
		return "", newErrorString(FileParseError, "the value of the basePath field is not string")
	}

	return "/", nil
}

func (p parserV2) getSchemes() ([]string, error) {
	schemesField, ok := p.data["schemes"]
	if !ok {
		switch p.url.Scheme {
		case "https", "http":
			return []string{p.url.Scheme}, nil
		default:
			return nil, newErrorString(FileParseError, "unsupported scheme")
		}
	}

	schemes, ok := schemesField.([]interface{})
	if !ok {
		return nil, newErrorString(FileParseError, "schemes field is not string array field")
	}
	if len(schemes) > 4 {
		return nil, newErrorString(FileParseError, "schemes length must be less than 4")
	}

	ret := make([]string, 0)
	for _, s := range supportedSchemes {
		if contains(schemes, s) {
			ret = append(ret, s)
		}
	}
	if len(ret) == 0 {
		return nil, newErrorString(FileParseError, "no supported schemes, i.e. https and http, is contained in schemes field")
	}
	return ret, nil
}

func (p parserV2) getBaseGatewayPath() (string, error) {
	basePathField, ok := p.data[basePathFieldName]
	if !ok {
		return "", newError(FileParseError, fmt.Errorf("%s is required field, but empty", basePathFieldName))
	}

	if basePath, ok := basePathField.(string); ok {
		if strings.HasPrefix(basePath, "/") {
			return basePath, nil
		}
		return "", newError(FileParseError, fmt.Errorf("%s field must start with '/'", basePathFieldName))
	}
	return "", newError(FileParseError, fmt.Errorf("%s field must be string field", basePathFieldName))
}

func (p parserV2) parsePaths() ([]API, error) {
	pathsField, ok := p.data["paths"]
	if !ok {
		return nil, newErrorString(FileParseError, "paths is required field, but empty")
	}

	paths, ok := pathsField.(map[string]interface{})
	if !ok {
		return nil, newErrorString(FileParseError, "paths field must be map field")
	}

	ret := make([]API, 0, len(paths))
	for path, value := range paths {
		if !strings.HasPrefix(path, "/") {
			return nil, newErrorString(FileParseError, "path's key must start with '/'")
		}
		description, ok := value.(map[string]interface{})
		if !ok {
			return nil, newError(FileParseError, fmt.Errorf("path's value must be map type, got non-map type in path %s", path))
		}

		var api API
		if xApidoorPathField, ok := description[pathFieldName]; ok {
			if xApidoorPath, ok := xApidoorPathField.(string); ok {
				if !strings.HasPrefix(xApidoorPath, "/") {
					return nil, newError(FileParseError, fmt.Errorf("%s start with '/', got %s in path %s", pathFieldName, xApidoorPath, path))
				}
				api = API{
					ForwardURL: path,
					Path:       xApidoorPath,
				}
			} else {
				return nil, newError(FileParseError, fmt.Errorf("%s must be string field, but %s in path %s is not string field", pathFieldName, pathFieldName, path))
			}
		} else {
			api = API{
				ForwardURL: path,
				Path:       path,
			}
		}
		ret = append(ret, api)
	}
	return ret, nil
}

// host field that contains scheme or path is invalid
func isOnlyHostContained(host string) bool {
	return !strings.ContainsRune(host, '/')
}

func contains(data []interface{}, target string) bool {
	for _, v := range data {
		str, ok := v.(string)
		if !ok {
			continue
		}
		if str == target {
			return true
		}
	}
	return false
}
