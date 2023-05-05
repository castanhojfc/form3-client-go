// Package form3 provides access to Form3 API using a client which as access to resources.
//
// Allows a client to perform http requests.
//
// Each resource has access to a set of operations.
package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultUrlScheme = "http"            // DefaultUrlScheme is the default URL scheme.
	DefaultUrlHost   = "accountapi:8080" // DefaultUrlHost is the default URL host.
	DefaultTimeout   = time.Second * 5   // DefaultTimeout is the default timeout on how much time should be used if no http response is obtained.
)

// Client is used to access API resourses.
type Client struct {
	BaseUrl    *url.URL      // API base Url to perform http requests.
	HttpClient *http.Client  // Http client used to perform http requests.
	Timeout    time.Duration // How much time should be used if no http response is obtained.

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
		Timeout:    DefaultTimeout,
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

func WithTimeout(timeout time.Duration) Option {
	return func(client *Client) {
		client.Timeout = timeout
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

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	request, error := http.NewRequest(method, requestURL, buffer)

	if error != nil {
		return nil, OperationError{Message: error.Error()}
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	request = request.WithContext(ctx)

	response, error := c.HttpClient.Do(request)

	if error != nil {
		return nil, OperationError{Message: error.Error()}
	}

	return response, nil
}
