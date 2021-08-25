package klient

import (
	"encoding/json"

	http "github.com/useflyent/fhttp"
)

type Response struct {
	Headers    http.Header
	Body       []byte
	Status     string
	StatusCode int
}

func (r *Response) JSON(data interface{}) error {
	return json.Unmarshal(r.Body, data)
}
