package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	oc_http "go.opencensus.io/plugin/ochttp"
	dd_http "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

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
		// Not all TokenProviders is thread-safe, wrap in CachedTokenProvider
		// to ensure it is. Is a no-op if it already is.
		c.TokenProvider = auth.NewCachedTokenProvider(provider)
	}
}

func WithDefaultHeader(header, value string) Option {
	return func(c *Client) {
		c.defaultHeaders.Set(header, value)
	}
}

func WithProblemDecoder(decoder ProblemDecoder) Option {
	return func(c *Client) {
		c.problemDecoder = decoder
	}
}

// WithOpenCensusTracing will add an OpenCensus transport to the client
// so that it will automatically inject trace-headers.
//
// Should be used when you trace your application with OpenCensus.
func WithOpenCensusTracing() Option {
	return func(c *Client) {
		c.client.Transport = new(oc_http.Transport)
	}
}

// WithTimeout sets http client timeout
//
// The default timeout of zero means no timeout.
//
// The Client cancels requests to the underlying Transport
// as if the Request's Context ended.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithDatadogTracing will add a Datadog transport to the client
// so that it will automatically inject trace-headers.
//
// Should be used when you trace your application with Datadog.
//
// Note, when used in AWS Lambda, make sure to set the service name.
// client.WithDatadogTracing(
//     dd_http.RTWithServiceName("<service_name>"),
// )
func WithDatadogTracing(opts ...dd_http.RoundTripperOption) Option {
	resourceNamer := func(req *http.Request) string {
		return fmt.Sprintf("%s %s", req.Method, req.URL.String())
	}
	opts = append([]dd_http.RoundTripperOption{
		dd_http.RTWithResourceNamer(resourceNamer),
	}, opts...)

	return func(c *Client) {
		c.client = dd_http.WrapClient(
			c.client,
			opts...,
		)
	}
}
