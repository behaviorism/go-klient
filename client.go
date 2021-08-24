package klient

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	client       *http.Client
	headersOrder *[]string
}

func NewClient(proxyURL string) *Client {
	var headersOrder []string

	transport := &transport{
		MaxIdleConns:        0,
		MaxConnsPerHost:     0,
		MaxIdleConnsPerHost: 100,
		headersOrder:        &headersOrder,
	}

	if len(proxyURL) > 0 {
		proxyUrl, _ := url.Parse(proxyURL)

		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	return &Client{
		client: &http.Client{
			Transport: transport,
		},
		headersOrder: &headersOrder,
	}
}

func (c *Client) SetHeadersOrder(headersList []string) {
	*c.headersOrder = headersList
}

func (c *Client) NewRequest() *Request {
	return &Request{
		client: c,
		header: make(http.Header),
	}
}

func (c *Client) AddCookie(u *url.URL, cookie *http.Cookie) error {
	if c.client.Jar == nil {
		return errors.New("Client doesn't have a cookie jar")
	}

	c.client.Jar.SetCookies(u, []*http.Cookie{cookie})

	return nil
}

func (c *Client) RemoveCookie(u *url.URL, cookie string) error {
	if c.client.Jar == nil {
		return errors.New("Client doesn't have a cookie jar")
	}

	newCookie := &http.Cookie{
		Name:  cookie,
		Value: "",
	}

	c.client.Jar.SetCookies(u, []*http.Cookie{newCookie})

	return nil
}

func (c *Client) Do(r *http.Request) (*Response, error) {
	resp, err := c.client.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	response := &Response{
		Headers:    resp.Header,
		Body:       body,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}

	return response, nil
}
