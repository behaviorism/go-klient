package klient

import (
	"errors"
	"io/ioutil"
	"net/url"

	http "github.com/useflyent/fhttp"
	"golang.org/x/net/proxy"
)

type Client struct {
	client *http.Client
}

type Browser struct {
	JA3       string
	UserAgent string
}

func NewClient(browser Browser, proxyURL string) *Client {
	if len(proxyURL) > 0 {
		dialer, _ := newConnectDialer(proxyURL)

		return &Client{
			client: &http.Client{Transport: newRoundTripper(browser, dialer)},
		}
	}

	return &Client{
		client: &http.Client{
			Transport: newRoundTripper(browser, proxy.Direct),
		},
	}
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
