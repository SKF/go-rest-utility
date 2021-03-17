package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
)

type RequestEcho struct {
	URL    string
	Method string
	Header http.Header
	Body   *string
}

func TestClientGet(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("endpoint")

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

func TestClientPut(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	payload := struct {
		ID     string
		Amount int
	}{
		"c8c9e607-219c-4b29-b161-474d4a331651", 2000,
	}

	request := Put("transfer/").WithJSONPayload(payload)

	client := NewClient(WithBaseURL(srv.URL))

	response, err := client.Do(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	echo := RequestEcho{}

	err = response.Unmarshal(&echo)
	require.NoError(t, err)
	require.Equal(t, "/transfer/", echo.URL)
	require.Equal(t, http.MethodPut, echo.Method)
	require.Equal(t, DefaultUserAgent, echo.Header.Get(headers.UserAgent))
	require.Equal(t, DefaultAcceptEncoding, echo.Header.Get(headers.AcceptEncoding))
	require.Equal(t, "application/json", echo.Header.Get(headers.ContentType))

	require.NotNil(t, echo.Body)
	require.Equal(t,
		`{"ID":"c8c9e607-219c-4b29-b161-474d4a331651","Amount":2000}`,
		strings.TrimSuffix(*echo.Body, "\n"),
	)
}

func TestClientDefaultHeader(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("endpoint")

	client := NewClient(
		WithBaseURL(srv.URL),
		WithDefaultHeader("X-Client-ID", "78147f11-62d9-4af0-917d-a0eb26d1c1fc"),
		WithDefaultHeader("User-Agent", "Custom"),
	)

	response, err := client.Do(context.Background(), request)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	echo := RequestEcho{}

	err = response.Unmarshal(&echo)
	require.NoError(t, err)
	require.Equal(t, "/endpoint", echo.URL)
	require.Equal(t, http.MethodGet, echo.Method)
	require.Equal(t, "Custom", echo.Header.Get(headers.UserAgent))
	require.Equal(t, DefaultAcceptEncoding, echo.Header.Get(headers.AcceptEncoding))
	require.Equal(t, "78147f11-62d9-4af0-917d-a0eb26d1c1fc", echo.Header.Get("X-Client-ID"))
}

// newEchoHTTPServer returns a new server which echos back the request as response.
func newEchoHTTPServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		var body *string
		if r.ContentLength != 0 {
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(rw, `{"error": "%s"}`, err)
			}

			str := string(bytes)
			body = &str
		}

		echo := RequestEcho{
			URL:    r.URL.String(),
			Method: r.Method,
			Header: r.Header,
			Body:   body,
		}

		if err := json.NewEncoder(rw).Encode(echo); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, `{"error": "%s"}`, err)
		}
	})

	return httptest.NewServer(handler)
}
