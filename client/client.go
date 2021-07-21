package client

import (
	"context"
	"fmt"
	"net/http"
	netUrl "net/url"

	"github.com/go-http-utils/headers"

	"github.com/SKF/go-rest-utility/client/auth"
)

const (
	DefaultUserAgent      string = "go-rest-utility/v1"
	DefaultAcceptEncoding string = "gzip"
)

type Client struct {
	BaseURL       *netUrl.URL
	TokenProvider auth.TokenProvider

	client         *http.Client
	defaultHeaders http.Header
}

// NewClient will create a new REST Client.
func NewClient(opts ...Option) *Client {
	client := &Client{
		BaseURL:        nil,
		TokenProvider:  nil,
		client:         new(http.Client),
		defaultHeaders: make(http.Header),
	}

	client.defaultHeaders.Set(headers.UserAgent, DefaultUserAgent)
	client.defaultHeaders.Set(headers.AcceptEncoding, DefaultAcceptEncoding)

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Do(ctx context.Context, r *Request) (*Response, error) {
	httpRequest, err := c.prepareRequest(ctx, r)
	if err != nil {
		return nil, err
	}

	httpResponse, err := c.client.Do(httpRequest)
	if err != nil {
		if httpResponse != nil && httpResponse.Body != nil {
			httpResponse.Body.Close()
		}

		return nil, fmt.Errorf("unable to perform http request: %w", err)
	}

	return &Response{*httpResponse}, nil
}

func (c *Client) DoAndUnmarshal(ctx context.Context, r *Request, v interface{}) error {
	response, err := c.Do(ctx, r)
	if err != nil {
		return err
	}

	return response.Unmarshal(v)
}

func (c *Client) prepareRequest(ctx context.Context, req *Request) (*http.Request, error) {
	url, err := req.ExpandURL(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid request URL: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, req.method, url.String(), req.body)
	if err != nil {
		return nil, fmt.Errorf("unable to create http request: %w", err)
	}

	for header, defaultValue := range c.defaultHeaders {
		if _, exists := req.header[header]; !exists {
			req.header[header] = defaultValue
		}
	}

	if c.TokenProvider != nil {
		token, err := c.TokenProvider.GetRawToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get token: %w", err)
		}

		req.header.Set(headers.Authorization, token.String())
	}

	httpRequest.Header = req.header

	return httpRequest, nil
}
