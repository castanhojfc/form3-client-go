//go:build unit

package form3_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/castanhojfc/form3-client-go/form3"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestForm3_New(t *testing.T) {
	t.Run("should create new client when no options are provided", func(t *testing.T) {
		t.Parallel()

		client, error := form3.New()

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, &url.URL{
			Scheme: form3.DefaultUrlScheme,
			Host:   form3.DefaultUrlHost,
		}, client.BaseUrl)
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

		httpClient := &http.Client{}
		client, error := form3.New(
			form3.WithHttpClient(&http.Client{}),
		)

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, &url.URL{
			Scheme: form3.DefaultUrlScheme,
			Host:   form3.DefaultUrlHost,
		}, client.BaseUrl)
		assert.Equal(t, httpClient, client.HttpClient)
	})

	t.Run("should not create client when base url option is malformed", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI("httsdf:fd8080")

		client, error := form3.New(
			form3.WithBaseUrl(url),
		)

		assert.Nil(t, client)
		expectedErrorMessage := "it was not possible to extract a scheme and a host from the provided URL"
		assert.Equal(t, form3.ClientError{Message: expectedErrorMessage}, error)
		assert.Equal(t, expectedErrorMessage, error.Error())
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
}

func TestForm3_PerformRequest(t *testing.T) {
	t.Run("should return an error when a malformed url is used", func(t *testing.T) {
		t.Parallel()

		client, _ := form3.New()
		response, error := client.PerformRequest("GET", "http://asdf.com/%%", []byte{})

		assert.Nil(t, response)
		assert.True(t, strings.Contains(error.Error(), "invalid URL escape"))
	})

	t.Run("should retry when service unavailable and stay within client http timeout", func(t *testing.T) {
		t.Parallel()
		// WARNING: This test is slow on purpose

		defer gock.Off()
		client, _ := form3.New(
			form3.WithHttpTimeout(2*time.Second),
			form3.WithHttpTimeUntilNextAttempt(1*time.Second),
			form3.WithDebugEnabled(false),
		)
		defer gock.RestoreClient(client.HttpClient)

		for i := 0; i <= 5; i++ {
			gock.New("http://test:8080").
				Get("/endpoint").
				Reply(503)
		}

		response, error := client.PerformRequest("GET", "http://test:8080/endpoint", []byte{})

		assert.Nil(t, response)
		assert.Equal(t, form3.OperationError{Message: "Get \"http://test:8080/endpoint\": context deadline exceeded", Body: nil}, error)
	})
}
