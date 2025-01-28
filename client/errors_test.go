package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client"
)

func TestClientGet_NotFoundError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "a nice description on why teapots are bad")
	}))
	defer srv.Close()

	request := client.Get("endpoint")

	c := client.NewClient(client.WithBaseURL(srv.URL))

	_, err := c.Do(context.Background(), request)
	require.Error(t, err)

	require.ErrorIs(t, err, client.ErrNotFound)

	httpErr := client.HTTPError{}
	require.ErrorAs(t, err, &httpErr)

	require.Equal(t, http.StatusNotFound, httpErr.StatusCode)
	require.Equal(t, "Not Found", httpErr.Status)
	require.Equal(t, "a nice description on why teapots are bad", httpErr.Body)
}
