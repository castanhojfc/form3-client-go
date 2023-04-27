package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	baseURL    string
	httpClient *http.Client

	Accounts *AccountService
}

func NewClient() *Client {
	client := &Client{
		baseURL:    os.Getenv("FORM3_ACCOUNT_API_URL"),
		httpClient: http.DefaultClient,
	}

	client.Accounts = &AccountService{client: client}

	return client
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	base, _ := url.Parse(c.baseURL)
	reference, _ := url.Parse(path)
	url := base.ResolveReference(reference)
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
