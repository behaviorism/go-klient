package klient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	client            *Client
	method, url, host string
	header            http.Header
	body              io.Reader
	err               error
}

func (r *Request) SetURL(url string) *Request {
	r.url = url

	return r
}

func (r *Request) SetMethod(method string) *Request {
	r.method = method

	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	if header, ok := r.header[key]; ok {
		header = append(header, value)
		r.header[key] = header
	} else {
		r.header[key] = []string{value}
	}

	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.header[key] = []string{value}

	return r
}

func (r *Request) SetHost(value string) *Request {
	r.host = value

	return r
}

func (r *Request) SetBodyString(body string) *Request {
	if r.err == nil {
		r.body = bytes.NewBufferString(body)
	}

	return r
}

func (r *Request) SetBodyJSON(body interface{}) *Request {
	b, err := json.Marshal(body)

	if err != nil {
		r.err = err
	} else {
		r.body = bytes.NewBuffer(b)
		r.header["content-type"] = []string{"application/json"}
	}

	return r
}

func (r *Request) SetBodyForm(body url.Values) *Request {
	if r.err == nil {
		r.body = strings.NewReader(body.Encode())
		r.header["content-type"] = []string{"application/x-www-form-urlencoded"}
	}

	return r
}

func (r *Request) Do() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	if len(r.url) == 0 {
		return nil, errors.New("URL is missing")
	}

	if len(r.method) == 0 {
		r.method = "GET"
	}

	req, err := http.NewRequest(r.method, r.url, r.body)

	if err != nil {
		return nil, err
	}

	req.Header = r.header

	if len(r.host) > 0 {
		req.Host = r.host
	}

	return r.client.Do(req)
}
