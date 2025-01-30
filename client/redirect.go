package client

import (
	"errors"
	"net/http"
)

type key int

var followRedirectsKey key

func redirectHandler(req *http.Request, via []*http.Request) error {
	// Default behavior from net/http
	if len(via) >= 10 { //nolint: mnd
		return errors.New("stopped after 10 redirects")
	}

	if follow, ok := req.Context().Value(followRedirectsKey).(bool); ok && !follow {
		return http.ErrUseLastResponse
	}

	return nil
}
