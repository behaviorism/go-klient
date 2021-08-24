package klient

import (
	"encoding/json"
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
				r.Header.Del(header)
			}
		}

		// merge remaining headers
		for header, value := range r.Header {
			newHeaders[header] = value
		}

		bs, _ := json.Marshal(newHeaders)

		fmt.Println("New headers", bs)

		r.Header = newHeaders
	}

	return t.roundTripper.RoundTrip(r)
}
