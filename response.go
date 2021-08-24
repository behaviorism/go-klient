package klient

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Headers    http.Header
	body       []byte
	Status     string
	StatusCode int
}

func (r *Response) JSON(data interface{}) error {
	return json.Unmarshal(r.body, data)
}

func (r *Response) Body() []byte {
	return r.body
}
