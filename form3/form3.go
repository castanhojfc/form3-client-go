package form3

import (
	"net/http"
	"os"
)

type Client struct {
	apiURL     string
	httpClient *http.Client
}

func NewClient() *Client {
	client := &Client{
		apiURL:     os.Getenv("FORM3_ACCOUNT_API_URL"),
		httpClient: http.DefaultClient,
	}

	return client
}
