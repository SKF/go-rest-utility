package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"

	. "github.com/SKF/go-rest-utility/client" //nolint: revive
)

func TestClientRedirects_AvoidFollowing(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("/redirect").
		WithFollowRedirects(false)

	client := NewClient(WithBaseURL(srv.URL))

	response, err := client.Do(context.Background(), request)
	require.NoError(t, err)

	require.Equal(t, http.StatusFound, response.StatusCode)
	require.Equal(t, "/", response.Header.Get("Location"))
}

func TestClientRedirects_Follows(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("/redirect{?to}").
		Assign("to", "/endpoint").
		WithFollowRedirects(true)

	client := NewClient(WithBaseURL(srv.URL))

	response, err := client.Do(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	echo := RequestEcho{}

	err = response.Unmarshal(&echo)
	require.NoError(t, err)
	require.Equal(t, "/endpoint", echo.URL)
	require.Equal(t, http.MethodGet, echo.Method)
	require.Equal(t, DefaultUserAgent, echo.Header.Get(headers.UserAgent))
	require.Equal(t, DefaultAcceptEncoding, echo.Header.Get(headers.AcceptEncoding))
}

func TestClientRedirects_Default(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("/redirect{?to}").
		Assign("to", "/endpoint")

	client := NewClient(WithBaseURL(srv.URL))

	response, err := client.Do(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	echo := RequestEcho{}

	err = response.Unmarshal(&echo)
	require.NoError(t, err)
	require.Equal(t, "/endpoint", echo.URL)
	require.Equal(t, http.MethodGet, echo.Method)
	require.Equal(t, DefaultUserAgent, echo.Header.Get(headers.UserAgent))
	require.Equal(t, DefaultAcceptEncoding, echo.Header.Get(headers.AcceptEncoding))
}
