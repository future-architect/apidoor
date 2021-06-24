package gateway

import "net/http"

func RequestHeaderSetter(src, dist *http.Request) {
	dist.Header = src.Header
	dist.Header.Del("Authorization")
	dist.Header.Del("Connection")
	dist.Header.Del("Cookie")
}
