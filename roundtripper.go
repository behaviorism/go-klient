package klient

import (
	"fmt"
	"net/http"
	"net/url"
)

type transport struct {
	headersOrder *[]string
	roundTripper http.RoundTripper

	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	Proxy               func(*http.Request) (*url.URL, error)
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println(r.Header)

	if r.Header != nil {
		var newHeaders http.Header

		for _, header := range *t.headersOrder {
			if value, ok := r.Header[header]; ok {
				newHeaders[header] = value
			}
		}

		r.Header = newHeaders
	}

	return t.roundTripper.RoundTrip(r)
}
