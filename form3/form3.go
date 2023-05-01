package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	reference, _ := url.Parse(path)
	url := c.baseUrl.ResolveReference(reference)
	url.Scheme = "http"
	buffer := newBuffer(body)
	request, error := http.NewRequest(
		method,
		url.String(),
		buffer,
	)

	if error != nil {
		return nil, error
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	request.Header.Set("Accept", "application/json")

	return request, nil
}

func newBuffer(body interface{}) io.ReadWriter {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	error := encoder.Encode(body)

	if error != nil {
		return nil
	}

	return buffer
}

func (c *Client) Do(ctx context.Context, request *http.Request, resource interface{}) (*http.Response, error) {
	response, error := c.httpClient.Do(request)

	if error != nil {
		return nil, error
	}

	defer response.Body.Close()

	if resource != nil {
		error = json.NewDecoder(response.Body).Decode(resource)
	}

	return response, error
}
