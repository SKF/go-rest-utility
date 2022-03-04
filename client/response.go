package client

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-http-utils/headers"
)

type Response struct {
	http.Response
}

func (r *Response) Unmarshal(v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("failed to json decode read bytes: %w", err)
	}

	return nil
}

type GzipReader struct {
	*gzip.Reader
	inner io.Closer
}

func (r *GzipReader) Close() error {
	// The underlying gzip.Reader assumes everything has been read for the checksum check to work.
	if _, err := io.Copy(io.Discard, r.Reader); err != nil {
		return fmt.Errorf(": %w", err)
	}

	if err := r.Reader.Close(); err != nil {
		return fmt.Errorf(": %w", err)
	}

	return r.inner.Close()
}

// DecompressResponse takes a http response and returns a decompressed
// http.Body and a set of headers that matches the decompressed result.
func DecompressResponse(resp http.Response) (io.ReadCloser, http.Header, error) {
	switch resp.Header.Get(headers.ContentEncoding) {
	case "gzip":
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return resp.Body, nil, err
		}

		resp.Header.Del(headers.ContentEncoding)

		return &GzipReader{
			Reader: gzipReader,
			inner:  resp.Body,
		}, resp.Header, nil
	default:
		return resp.Body, resp.Header, nil
	}
}
