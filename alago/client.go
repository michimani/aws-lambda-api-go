package alago

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	runtimeAPIEnvKey string = "AWS_LAMBDA_RUNTIME_API"
)

var (
	defaultHttpClient *http.Client = http.DefaultClient
)

type AlagoClient interface {
	Host() string
	HttpClient() *http.Client
}

// Client is the struct for call runtime api.
type Client struct {
	// http.Client for call runtime api
	c *http.Client

	// host of runtime api endpoint
	host string
}

// NewClientInput is the struct for creating new Client.
// If HttpClient is nil, default http.Client will be used for Client.
type NewClientInput struct {
	HttpClient *http.Client
}

// NewClient returns new Client.
// If environment variable `AWS_LAMBDA_RUNTIME_API` is not set or
// the value of it is empty, returns a error.
func NewClient(in *NewClientInput) (*Client, error) {
	if in == nil {
		return nil, errors.New("NewClientInput is nil.")
	}

	host := os.Getenv(runtimeAPIEnvKey)
	if host == "" {
		return nil, fmt.Errorf("%s is not set or the value of it is empty.", runtimeAPIEnvKey)
	}

	var c *http.Client
	if in.HttpClient != nil {
		c = in.HttpClient
	} else {
		c = defaultHttpClient
	}

	return &Client{
		c:    c,
		host: host,
	}, nil
}

func (c *Client) Host() string {
	if c == nil {
		return ""
	}
	return c.host
}

func (c *Client) HttpClient() *http.Client {
	if c == nil {
		return nil
	}
	return c.c
}
