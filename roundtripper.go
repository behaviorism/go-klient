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
	if len(r.Header) > 0 {
		newHeaders := make(http.Header)

		// add headers within order with priority
		for _, header := range *t.headersOrder {
			if value, ok := r.Header[header]; ok {
				newHeaders[header] = value
				fmt.Println("Header", newHeaders)
				r.Header.Del(header)
			}
		}

		// merge remaining headers
		for header, value := range r.Header {
			newHeaders[header] = value
		}

		r.Header = newHeaders
	}

	return t.roundTripper.RoundTrip(r)
}
