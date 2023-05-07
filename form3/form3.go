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
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultUrlScheme                = "http"            // DefaultUrlScheme is the default URL scheme.
	DefaultUrlHost                  = "accountapi:8080" // DefaultUrlHost is the default URL host.
	DefaultHttpTimeout              = time.Second * 60  // DefaultTimeout is the default timeout on how much time should be used if no http response is obtained.
	DefaultHttpRetryAttempts        = 3                 // DefaultHttpRetryAttempts is the default number of attempts when performing a http request.
	DefaultHttpTimeUntilNextAttempt = 1 * time.Second   // DefaultHttpTimeUntilNextAttempt is the default time until the next http request attemp is made.
	DefaultDebugEnabled             = false             // DefaultDebug is the default value to determine if debug messages shall be shown.
)

// LogDebugMessage defines the function interface that is used to log debug messages.
type LogDebugMessage func(format string, v ...any)

// Client is used to access API resourses.
type Client struct {
	BaseUrl                   *url.URL        // API base Url to perform http requests.
	HttpClient                *http.Client    // Http client used to perform http requests.
	HttpTimeout               time.Duration   // How much time should be used if no http response is obtained.
	HttpRetryAttempts         int             // How many attempts shall be made if an http cannot be made but can be retried.
	HttpTimeUntilNextAttempt  time.Duration   // How much time should be spent until the next http retry attempt is done.
	DebugEnabled              bool            // If debugging messages should be shown.
	HttpRetryJitterRandomSeed rand.Source     // Random seed used to generate jitter between http retry attempts.
	Accounts                  *AccountService // Account Service, has access to operations.
	UserAgent                 string          // Allow the server to identify the client.
	LogDebugMessage           LogDebugMessage // Allow the client to log debug messages.
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
		HttpClient:                http.DefaultClient,
		HttpTimeout:               DefaultHttpTimeout,
		HttpRetryAttempts:         DefaultHttpRetryAttempts,
		HttpTimeUntilNextAttempt:  DefaultHttpTimeUntilNextAttempt,
		DebugEnabled:              DefaultDebugEnabled,
		HttpRetryJitterRandomSeed: rand.NewSource(time.Now().UnixNano()),
		UserAgent:                 "form3-client-go",
	}

	client.Accounts = &AccountService{Client: client, JsonMarshal: json.Marshal, JsonUnmarshal: json.Unmarshal, ReadAll: io.ReadAll}
	client.LogDebugMessage = log.Printf

	return client, nil
}

// PerformRequest uses a client to perform a http request to the API.
//
// An error is returned if there was any problem creating or performing the request.
// Requests can be retried if possible. The time until the next attempt is doubled but it stays within the http timeout.
// Some jitter is added between requests.
func (c *Client) PerformRequest(method string, requestURL string, body []byte) (*http.Response, error) {
	var buffer io.ReadWriter

	if body != nil {
		buffer = bytes.NewBuffer(body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.HttpTimeout)
	defer cancel()

	request, _error := http.NewRequest(method, requestURL, buffer)

	if _error != nil {
		return nil, OperationError{Message: _error.Error()}
	}

	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	request.Header.Set("User-Agent", c.UserAgent)

	request = request.WithContext(ctx)

	response, _error := c.retryRequest(c.HttpRetryAttempts, c.HttpTimeUntilNextAttempt, func() (*http.Response, error) {
		return c.HttpClient.Do(request)
	})

	if _error != nil {
		return nil, OperationError{Message: _error.Error()}
	}

	return response, nil
}

func (c *Client) retryRequest(remainingAttempts int, timeUntilNextAttempt time.Duration, retriable func() (*http.Response, error)) (*http.Response, error) {
	response, error := retriable()

	if response == nil {
		return nil, error
	}

	// Do not retry on client errors. If the client performed too many requests it is still possible to retry.
	if response.StatusCode >= 400 && response.StatusCode < 500 && response.StatusCode != 429 {
		return response, error
	}

	if error != nil || response.StatusCode >= 500 || response.StatusCode == 429 {
		if remainingAttempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(timeUntilNextAttempt))) / 3
			timeUntilNextAttempt = (timeUntilNextAttempt * 2) + jitter

			// Keep the next attempt within the client timeout
			if timeUntilNextAttempt > c.HttpTimeout {
				timeUntilNextAttempt = c.HttpTimeout
			}

			if c.DebugEnabled {
				c.LogDebugMessage("DEBUG: Http request failed, retrying in: %v jitter addded: %v remaining attempts: %d", timeUntilNextAttempt, jitter, remainingAttempts)
			}

			remainingAttempts--
			time.Sleep(timeUntilNextAttempt)

			return c.retryRequest(remainingAttempts, timeUntilNextAttempt, retriable)
		}
	}

	return response, error
}
