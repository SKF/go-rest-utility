package client

import (
	"net/url"

	"github.com/SKF/go-rest-utility/client/auth"
)

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		// If the provided URL is not valid the final URL will not be valid and
		// therefore it is safe to ignore this error.
		c.BaseURL, _ = url.Parse(baseURL) //nolint:errcheck
	}
}

func WithTokenProvider(provider auth.TokenProvider) Option {
	return func(c *Client) {
		c.TokenProvider = provider
	}
}
