// Package form3 provides access to Form3 API using a client which as access to resources.
//
// Allows a client to perform http requests.
//
// Each resource has access to a set of operations.
package form3

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	DefaultUrlScheme = "http"            // DefaultUrlScheme is the default URL scheme
	DefaultUrlHost   = "accountapi:8080" // DefaultUrlHost is the default URL host
)

// Client is used to access API resourses.
type Client struct {
	BaseUrl    *url.URL     // API base URL, can be configured.
	HttpClient *http.Client // Http client, can be configured.

	Accounts *AccountService // Account Service, has access to operations.
}

// Option represents an option that can be externally configured.
type Option func(c *Client)

// New creates a new client.
//
// A set of options can be used to customize it.
//
// An error is returned if there is a problem finding a URL scheme and host.
func New(options ...Option) (*Client, error) {
	client := &Client{
		BaseUrl: &url.URL{
			Scheme: DefaultUrlScheme,
			Host:   DefaultUrlHost,
		},
		HttpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	if client.BaseUrl.Scheme == "" || client.BaseUrl.Host == "" {
		return nil, ClientError{
			Message: "it was not possible to extract a scheme and a host from the provided URL",
		}
	}

	client.Accounts = &AccountService{Client: client, JsonMarshal: json.Marshal, JsonUnmarshal: json.Unmarshal, ReadAll: io.ReadAll}

	return client, nil
}

// WithBaseUrl allows the base URL to be configured externally.
func WithBaseUrl(url *url.URL) Option {
	return func(client *Client) {
		client.BaseUrl = url
	}
}

// WithHttpClient allows the http client to be configured externally.
func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.HttpClient = httpClient
	}
}

// PerformRequest uses a client to perform a http request to the API.
//
// An error is returned if there was any problem creating or performing the request.
func PerformRequest(c *Client, method string, requestURL string, body []byte) (*http.Response, error) {
	var buffer io.ReadWriter

	if body != nil {
		buffer = bytes.NewBuffer(body)
	}

	request, error := http.NewRequest(method, requestURL, buffer)

	if error != nil {
		return nil, OperationError{Message: error.Error()}
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, error := c.HttpClient.Do(request)

	if error != nil {
		return nil, OperationError{Message: error.Error()}
	}

	return response, nil
}
