package form3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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

func performRequest(c *Client, method string, requestURL string, body []byte) (*http.Response, error) {
	var buffer io.ReadWriter

	if body != nil {
		buffer = bytes.NewBuffer(body)
	}

	request, error := http.NewRequest(method, requestURL, buffer)

	if error != nil {
		return nil, fmt.Errorf("there was a problem creating the request: %w", error)
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, error := c.httpClient.Do(request)

	if error != nil {
		return nil, fmt.Errorf("there was a problem performing the request: %w", error)
	}

	return response, nil
}

func buildUnsuccessfulResponse(response *http.Response) error {
	dump, error := httputil.DumpResponse(response, true)

	if error != nil {
		return fmt.Errorf("it was not possible dump the response for an unsucessful operation: %w", error)
	}

	return OperationError{
		Message:  "could not perform operation",
		Response: dump,
	}
}
