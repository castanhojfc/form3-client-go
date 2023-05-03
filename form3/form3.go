package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Client struct {
	BaseUrl    *url.URL
	HttpClient *http.Client

	Accounts *AccountService
}

type Option func(c *Client)

func New(options ...Option) (*Client, error) {
	client := &Client{
		BaseUrl:    nil,
		HttpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	if client.BaseUrl == nil {
		rawUrl := os.Getenv("FORM3_ACCOUNT_API_URL")

		if rawUrl == "" {
			return nil, fmt.Errorf("no base URL was provided, it should be set using environment variable or as an option")
		}

		parsedUrl, error := url.ParseRequestURI(rawUrl)

		if error != nil {
			return nil, fmt.Errorf("there was a problem parsing the URL: %w", error)
		}

		client.BaseUrl = parsedUrl
	}

	if client.BaseUrl.Scheme == "" || client.BaseUrl.Host == "" {
		return nil, fmt.Errorf("it was not possible to extract a scheme and a host from the provided URL: %s", client.BaseUrl)
	}

	client.Accounts = &AccountService{Client: client, JsonMarshal: json.Marshal, JsonUnmarshal: json.Unmarshal, ReadAll: io.ReadAll}

	return client, nil
}

func WithBaseUrl(url *url.URL) Option {
	return func(client *Client) {
		client.BaseUrl = url
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.HttpClient = httpClient
	}
}

func PerformRequest(c *Client, method string, requestURL string, body []byte) (*http.Response, error) {
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

	response, error := c.HttpClient.Do(request)

	if error != nil {
		return nil, fmt.Errorf("there was a problem performing the request: %w", error)
	}

	return response, nil
}

func BuildUnsuccessfulResponse(response *http.Response) error {
	dump, _ := httputil.DumpResponse(response, true)

	return OperationError{
		Message:  "could not perform operation",
		Response: string(dump),
	}
}
