//go:build unit

package form3_test

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/castanhojfc/form3-client-go/form3"
	"github.com/stretchr/testify/assert"
)

func TestForm3_New(t *testing.T) {
	t.Run("should create new client when no options are provided", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI(os.Getenv("FORM3_ACCOUNT_API_URL"))
		client, error := form3.New()

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, url, client.BaseUrl)
		assert.Equal(t, http.DefaultClient, client.HttpClient)
	})

	t.Run("should create new client with url when option is provided", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI("http://asdf:8080")
		client, error := form3.New(
			form3.WithBaseUrl(url),
		)

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, url, client.BaseUrl)
		assert.Equal(t, http.DefaultClient, client.HttpClient)
	})

	t.Run("should create new client with http client when option is provided", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI(os.Getenv("FORM3_ACCOUNT_API_URL"))
		httpClient := &http.Client{}
		client, error := form3.New(
			form3.WithHttpClient(&http.Client{}),
		)

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, url, client.BaseUrl)
		assert.Equal(t, httpClient, client.HttpClient)
	})

	t.Run("should not create client when base url option is malformed, gives priority to option", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI("httsdf:fd8080")

		client, error := form3.New(
			form3.WithBaseUrl(url),
		)

		assert.Nil(t, client)
		assert.Equal(t, "it was not possible to extract a scheme and a host from the provided URL: httsdf:fd8080", error.Error())
	})

	t.Run("should create new client with all options provided", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI("http://asdf:8080")
		httpClient := &http.Client{}
		client, error := form3.New(
			form3.WithBaseUrl(url),
			form3.WithHttpClient(&http.Client{}),
		)

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, url, client.BaseUrl)
		assert.Equal(t, httpClient, client.HttpClient)
	})

	t.Run("should not create client when no base url is provided", func(t *testing.T) {
		// WARNING: This test cannot be run in parallel
		t.Setenv("FORM3_ACCOUNT_API_URL", "")

		client, error := form3.New()

		assert.Nil(t, client)
		assert.Equal(t, "no base URL was provided, it should be set using environment variable or as an option", error.Error())
	})

	t.Run("should not create client when base url by environment variable is malformed", func(t *testing.T) {
		// WARNING: This test cannot be run in parallel
		t.Setenv("FORM3_ACCOUNT_API_URL", "fasdf:3030")

		client, error := form3.New()

		assert.Nil(t, client)
		assert.Equal(t, "it was not possible to extract a scheme and a host from the provided URL: fasdf:3030", error.Error())
	})

	t.Run("should not create client when base url by environment variable and option are malformed", func(t *testing.T) {
		// WARNING: This test cannot be run in parallel
		t.Setenv("FORM3_ACCOUNT_API_URL", "fasdf:3030")
		url, _ := url.ParseRequestURI("httsdf:fd8080")

		client, error := form3.New(
			form3.WithBaseUrl(url),
		)

		assert.Nil(t, client)
		assert.Equal(t, "it was not possible to extract a scheme and a host from the provided URL: httsdf:fd8080", error.Error())
	})

	t.Run("should create client when no base url is provided through environment variable but option is provided", func(t *testing.T) {
		// WARNING: This test cannot be run in parallel
		t.Setenv("FORM3_ACCOUNT_API_URL", "")
		url, _ := url.ParseRequestURI("http://asdf:8080")

		client, error := form3.New(
			form3.WithBaseUrl(url),
		)

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, url, client.BaseUrl)
	})

	t.Run("should not create client when there is a problem parsing the url via environment variable", func(t *testing.T) {
		// WARNING: This test cannot be run in parallel
		t.Setenv("FORM3_ACCOUNT_API_URL", "---")

		client, error := form3.New()

		assert.Nil(t, client)
		assert.Equal(t, "there was a problem parsing the URL: parse \"---\": invalid URI for request", error.Error())
	})
}

func TestForm3_PerformRequest(t *testing.T) {
	t.Run("should return an error when a malformed url is used", func(t *testing.T) {
		client, _ := form3.New()
		response, error := form3.PerformRequest(client, "GET", "http://asdf.com/%%", []byte{})

		assert.Nil(t, response)
		assert.True(t, strings.Contains(error.Error(), "invalid URL escape"))
	})
}
