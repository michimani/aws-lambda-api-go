package alago

import "net/http"

func NewClientWithClientAndHost(c *http.Client, h string) *Client {
	return &Client{
		host: h,
		c:    c,
	}
}
