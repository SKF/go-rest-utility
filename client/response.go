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
	defer r.Close()

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("failed to json decode read bytes: %w", err)
	}

	return nil
}

// Close reads all of the body stream and closes it to make sure that tcp connections can be reused properly
func (r *Response) Close() error {
	if _, err := io.Copy(io.Discard, r.Body); err != nil {
		r.Body.Close() // nolint: errcheck

		return err
	}

	return r.Body.Close()
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
// If the content-length header is 0, return the body and the header
// without decompressing.
func DecompressResponse(resp http.Response) (io.ReadCloser, http.Header, error) {
	if contentLengthHeader := resp.Header.Get(headers.ContentLength); contentLengthHeader == "0" {
		return resp.Body, resp.Header, nil
	}

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
