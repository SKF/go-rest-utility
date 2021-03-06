package client

import (
	"compress/gzip"
	"encoding/json"
	"net/http"

	"github.com/go-http-utils/headers"
)

type Response struct {
	http.Response
}

func (r *Response) Unmarshal(v interface{}) (err error) {
	reader := r.Body

	switch r.Header.Get(headers.ContentEncoding) {
	case "gzip":
		defer reader.Close()

		if reader, err = gzip.NewReader(reader); err != nil {
			return
		}
	}

	defer reader.Close()

	if err = json.NewDecoder(reader).Decode(&v); err != nil {
		return
	}

	return nil
}
