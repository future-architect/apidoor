package swagger_parser

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type parserV3 struct {
	*Parser
}

func newParserV3(base *Parser) *parserV3 {
	return &parserV3{
		base,
	}
}

func (p parserV3) parse() (*Swagger, error) {
	serverURL, err := p.getForwardURLBaseAndScheme()
	if err != nil {
		return nil, err
	}
	forwardURLBase := ""
	if serverURL.Port() == "" {
		forwardURLBase = fmt.Sprintf("%s%s", serverURL.Host, serverURL.Path)
	}
	schemes := []string{serverURL.Scheme}

	pathBase, err := p.getBaseGatewayPath()
	if err != nil {
		return nil, err
	}

	apis, err := p.parsePaths()
	if err != nil {
		return nil, err
	}

	return &Swagger{
		Version:        "v3",
		ForwardURLBase: forwardURLBase,
		Schemes:        schemes,
		PathBase:       pathBase,
		APIs:           apis,
	}, nil
}

func (p parserV3) getForwardURLBaseAndScheme() (*url.URL, error) {
	serversField, ok := p.data["servers"]
	if !ok {
		serverURL := p.url
		serverURL.Path = filepath.Dir(serverURL.Path)
		return serverURL, nil
	}
	servers, ok := serversField.([]interface{})
	if !ok {
		return nil, newErrorString(FileParseError, "servers field must be array of maps")
	}

	urls, err := p.parseForwardURLFromServers(servers)
	if err != nil {
		return nil, err
	}

	if len(urls) == 0 {
		return nil, newErrorString(FileParseError, "no valid forward url provided")
	}
	// apidoor do not support multiple forward URL
	return urls[0], nil
}

func (p parserV3) parseForwardURLFromServers(servers []interface{}) ([]*url.URL, error) {

	ret := make([]*url.URL, 0, len(servers))
	for _, server := range servers {
		serverMap, ok := server.(map[string]interface{})
		if !ok {
			continue
		}
		urlField, ok := serverMap["url"]
		if !ok {
			return nil, newErrorString(FileParseError, "url is a require field in each server object")
		}
		urlStr, ok := urlField.(string)
		if !ok {
			return nil, newErrorString(FileParseError, "url must be a string field")
		}
		serverUrl, err := url.Parse(urlStr)
		if err != nil {
			return nil, newErrorString(FileParseError, "url field value is not url format")
		}

		// apidoor does not support URL template in base path
		if _, ok := serverMap["variables"]; ok {
			return nil, newErrorString(FileParseError, "apidoor does not support URL template in base path")
		}

		// url is relative path
		if serverUrl.Host == "" {
			documentDir := filepath.Dir(p.url.Path)
			path := filepath.Clean(documentDir + serverUrl.Path)
			serverUrl = p.url
			serverUrl.Path = path
		}

		switch serverUrl.Scheme {
		case "https", "http":
			ret = append(ret, serverUrl)
		}
	}
	return ret, nil
}

func (p parserV3) getBaseGatewayPath() (string, error) {
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

func (p parserV3) parsePaths() ([]API, error) {
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
