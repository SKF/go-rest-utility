package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
)

func TestClientGet(t *testing.T) {
	srv := newEchoHTTPServer()
	defer srv.Close()

	request := Get("endpoint")

	client := NewClient(WithBaseURL(srv.URL))

	response, err := client.Do(context.Background(), request)

	if assert.NoError(t, err) {
		assert.Equal(t, 200, response.StatusCode)

		incomingRequest := IncomingRequest{}

		err := response.Unmarshal(&incomingRequest)
		if assert.NoError(t, err) {
			assert.Equal(t, "/endpoint", incomingRequest.URL)
			assert.Equal(t, DefaultUserAgent, incomingRequest.Header.Get(headers.UserAgent))
			assert.Equal(t, DefaultAcceptEncoding, incomingRequest.Header.Get(headers.AcceptEncoding))
		}
	}
}

type IncomingRequest struct {
	URL    string
	Header http.Header
}

func newEchoHTTPServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		incomingRequest := IncomingRequest{
			URL:    r.URL.String(),
			Header: r.Header,
		}

		if err := json.NewEncoder(rw).Encode(incomingRequest); err != nil {
			rw.WriteHeader(500)
			fmt.Fprintf(rw, `{"error": "%s"}`, err)
		}
	})

	return httptest.NewServer(handler)
}
