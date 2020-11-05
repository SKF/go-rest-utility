package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-http-utils/headers"
	"go.opencensus.io/plugin/ochttp"

	"github.com/SKF/go-rest-utility/client/auth"
)

const (
	DefaultUserAgent      string = "go-rest-utility/v1"
	DefaultAcceptEncoding string = "gzip"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL       *url.URL
	TokenProvider auth.TokenProvider

	client Doer
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		BaseURL:       nil,
		TokenProvider: nil,
		client: &http.Client{
			Transport: &ochttp.Transport{},
		},
	}

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
		return nil, fmt.Errorf("unable to perform http request: %w", err)
	}

	return c.prepareResponse(httpResponse)
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

	if req.header.Get(headers.UserAgent) == "" {
		req.header.Set(headers.UserAgent, DefaultUserAgent)
	}

	if c.TokenProvider != nil {
		token, err := c.TokenProvider.GetRawToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get token: %w", err)
		}

		req.header.Set(headers.Authorization, token.String())
	}

	req.header.Set(headers.AcceptEncoding, DefaultAcceptEncoding)

	httpRequest.Header = req.header

	return httpRequest, nil
}

func (c *Client) prepareResponse(resp *http.Response) (*Response, error) {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()

		errorBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("failed to get: %s, got status code: %d", resp.Request.URL, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("got 401 for %s: %s: %w", resp.Request.URL, errorBody, ErrUnauthorized)
		} else if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("got 403 for %s: %s: %w", resp.Request.URL, errorBody, ErrForbidden)
		} else if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("got 404 for %s: %s: %w", resp.Request.URL, errorBody, ErrNotFound)
		}

		return nil, fmt.Errorf("failed to get: %s, got status code: %d, body: %s", resp.Request.URL, resp.StatusCode, errorBody)
	}

	return &Response{*resp}, nil
}
