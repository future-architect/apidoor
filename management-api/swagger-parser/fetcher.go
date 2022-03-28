package swaggerparser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type contentType int

const (
	typeJson contentType = iota
	typeYaml
)

var (
	NotFoundErr error = errors.New("file not found")
)

type swaggerFile struct {
	body        *bytes.Buffer
	contentType contentType
}

type FileFetcher interface {
	Fetch(ctx context.Context, url *url.URL) (*swaggerFile, error)
}

type DefaultFetcher struct {
	client *http.Client
}

func NewDefaultFetcher() DefaultFetcher {
	return DefaultFetcher{
		client: &http.Client{},
	}
}

func (df DefaultFetcher) Fetch(ctx context.Context, url *url.URL) (*swaggerFile, error) {
	// TODO: retry etc.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, newError(FetchClientError, fmt.Errorf("create http request: %w", err))
	}

	resp, err := df.client.Do(req)
	if err != nil {
		return nil, newError(FetchServerError, fmt.Errorf("fetch file: %w", err))
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	if _, err = io.Copy(body, resp.Body); err != nil {
		return nil, newError(FetchClientError, fmt.Errorf("read response body: %w", err))
	}

	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, newError(FormatError, NotFoundErr)
		default:
			return nil, newError(FetchServerError, fmt.Errorf("http status is not 200, status %d", resp.StatusCode))
		}
	}

	var contentType contentType
	header := resp.Header.Get("Content-Type")
	switch header {
	case "application/json":
		contentType = typeJson
	case "text/x-yaml", "application/x-yaml", "text/vnd.yaml":
		contentType = typeYaml
	default:
		return nil, newError(FetchServerError, fmt.Errorf("content-type is not set or unsupported, content-type %s", header))
	}

	return &swaggerFile{
		body:        body,
		contentType: contentType,
	}, nil
}
