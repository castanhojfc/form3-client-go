package form3

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	baseUrl    *url.URL
	httpClient *http.Client

	Accounts *AccountService
}

type Option func(c *Client)

func NewClient(opts ...Option) (*Client, error) {
	rawUrl := os.Getenv("FORM3_ACCOUNT_API_URL")
	baseUrl, error := url.ParseRequestURI(rawUrl)

	if error != nil {
		return nil, fmt.Errorf("there was a problem parsing the URL: %w", error)
	}

	client := &Client{
		baseUrl:    baseUrl,
		httpClient: http.DefaultClient,
	}

	for _, o := range opts {
		o(client)
	}

	client.Accounts = &AccountService{client: client}

	return client, nil
}

func WithBaseUrl(url *url.URL) Option {
	return func(c *Client) {
		c.baseUrl = url
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
