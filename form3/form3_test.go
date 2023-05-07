//go:build unit

package form3_test

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/castanhojfc/form3-client-go/form3"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		client, error := form3.New()
		client.BaseUrl = url

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, url, client.BaseUrl)
		assert.Equal(t, http.DefaultClient, client.HttpClient)
	})

	t.Run("should create new client with http client when option is provided", func(t *testing.T) {
		t.Parallel()

		httpClient := &http.Client{}
		client, error := form3.New()
		client.HttpClient = &http.Client{}

		assert.NotNil(t, client)
		assert.Nil(t, error)
		assert.Equal(t, &url.URL{
			Scheme: form3.DefaultUrlScheme,
			Host:   form3.DefaultUrlHost,
		}, client.BaseUrl)
		assert.Equal(t, httpClient, client.HttpClient)
	})

	t.Run("should create new client with all options provided", func(t *testing.T) {
		t.Parallel()

		url, _ := url.ParseRequestURI("http://asdf:8080")
		httpClient := &http.Client{}
		client, error := form3.New()
		client.BaseUrl = url
		client.HttpClient = &http.Client{}
		client.DebugEnabled = true
		client.HttpRetryAttempts = 4
		client.HttpTimeUntilNextAttempt = 3 * time.Second
		client.HttpTimeout = 10 * time.Second

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

	t.Run("should retry when service unavailable and return the successful response after retries", func(t *testing.T) {
		defer gock.Off()
		client, _ := form3.New()
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = 100

		for i := 0; i <= 2; i++ {
			gock.New("http://test:8080").
				Get("/endpoint").
				Reply(503)
		}

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(200).
			JSON(map[string]string{"outcome": "success"})

		response, error := client.PerformRequest("GET", "http://test:8080/endpoint", nil)

		assert.Nil(t, error)
		assert.Equal(t, 200, response.StatusCode)
		body, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, []byte("{\"outcome\":\"success\"}\n"), body)
	})

	t.Run("do not retry when number of retry attempts is configured to be less than 1", func(t *testing.T) {
		defer gock.Off()
		client, _ := form3.New()
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = -666

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(503)

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(200).
			JSON(map[string]string{"outcome": "success"})

		response, error := client.PerformRequest("GET", "http://test:8080/endpoint", nil)

		assert.Nil(t, error)
		assert.Equal(t, 503, response.StatusCode)
	})

	t.Run("do not retry when the first response contains a client error status code", func(t *testing.T) {
		defer gock.Off()
		client, _ := form3.New()
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = 100

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(400)

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(200).
			JSON(map[string]string{"outcome": "success"})

		response, error := client.PerformRequest("GET", "http://test:8080/endpoint", nil)

		assert.Nil(t, error)
		assert.Equal(t, 400, response.StatusCode)
	})

	t.Run("retry when the response contains too many requests client error status code", func(t *testing.T) {
		defer gock.Off()
		client, _ := form3.New()
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = 100

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(429)

		gock.New("http://test:8080").
			Get("/endpoint").
			Reply(200).
			JSON(map[string]string{"outcome": "success"})

		response, error := client.PerformRequest("GET", "http://test:8080/endpoint", nil)

		assert.Nil(t, error)
		assert.Equal(t, 200, response.StatusCode)
		body, _ := ioutil.ReadAll(response.Body)
		assert.Equal(t, []byte("{\"outcome\":\"success\"}\n"), body)
	})

	t.Run("should retry when service unavailable and stay within client http timeout", func(t *testing.T) {
		defer gock.Off()
		client, _ := form3.New()
		client.HttpTimeout = 100 * time.Microsecond
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = 100

		for i := 0; i <= 5; i++ {
			gock.New("http://test:8080").
				Get("/endpoint").
				Reply(503)
		}

		response, error := client.PerformRequest("GET", "http://test:8080/endpoint", []byte{})

		assert.Nil(t, response)
		assert.Equal(t, form3.OperationError{Message: "Get \"http://test:8080/endpoint\": context deadline exceeded", Body: nil}, error)
	})

	t.Run("should retry when service unavailable and print debug messages if debug is enabled", func(t *testing.T) {
		defer gock.Off()
		mockLogDebugMessage := new(LogDebugMessageMock)
		mockLogDebugMessage.On("LogDebugMessage", mock.Anything, mock.Anything).Return()

		client, _ := form3.New()
		client.HttpTimeout = 100 * time.Microsecond
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = 1
		client.DebugEnabled = true
		rand.Seed(0)

		client.LogDebugMessage = mockLogDebugMessage.LogDebugMessage

		for i := 0; i <= 3; i++ {
			gock.New("http://test:8080").
				Get("/endpoint").
				Reply(503)
		}

		client.PerformRequest("GET", "http://test:8080/endpoint", []byte{})

		mockLogDebugMessage.AssertCalled(t, "LogDebugMessage", "DEBUG: Http request failed, retrying in: %v jitter addded: %v remaining attempts: %d", []interface{}{time.Duration(100000), time.Duration(5168), 1})
	})

	t.Run("should retry when service unavailable and not print debug messages if debug is disabled", func(t *testing.T) {
		defer gock.Off()
		mockLogDebugMessage := new(LogDebugMessageMock)
		mockLogDebugMessage.On("LogDebugMessage", mock.Anything, mock.Anything).Return()

		client, _ := form3.New()
		client.HttpTimeout = 100 * time.Microsecond
		client.HttpTimeUntilNextAttempt = 50 * time.Microsecond
		client.HttpRetryAttempts = 1
		client.DebugEnabled = false
		rand.Seed(0)

		client.LogDebugMessage = mockLogDebugMessage.LogDebugMessage

		for i := 0; i <= 3; i++ {
			gock.New("http://test:8080").
				Get("/endpoint").
				Reply(503)
		}

		client.PerformRequest("GET", "http://test:8080/endpoint", []byte{})

		mockLogDebugMessage.AssertNotCalled(t, "LogDebugMessage")
	})
}
