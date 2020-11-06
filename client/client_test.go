package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
)

func TestClientGet(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("endpoint")

	client := NewClient(WithBaseURL(srv.URL))

	response, err := client.Do(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)

	echo := RequestEcho{}

	err = response.Unmarshal(&echo)
	require.NoError(t, err)
	require.Equal(t, "/endpoint", echo.URL)
	require.Equal(t, http.MethodGet, echo.Method)
	require.Equal(t, DefaultUserAgent, echo.Header.Get(headers.UserAgent))
	require.Equal(t, DefaultAcceptEncoding, echo.Header.Get(headers.AcceptEncoding))
}

type RequestEcho struct {
	URL    string
	Method string
	Header http.Header
}

// newEchoHTTPServer returns a new server which echos back the request as response.
func newEchoHTTPServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		echo := RequestEcho{
			URL:    r.URL.String(),
			Method: r.Method,
			Header: r.Header,
		}

		if err := json.NewEncoder(rw).Encode(echo); err != nil {
			rw.WriteHeader(500)
			fmt.Fprintf(rw, `{"error": "%s"}`, err)
		}
	})

	return httptest.NewServer(handler)
}
