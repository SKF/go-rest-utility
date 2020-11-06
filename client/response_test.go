package client

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
)

func TestResponseUnmarshalSimple(t *testing.T) {
	response := Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			Header:     make(http.Header),
		},
	}

	value := struct {
		Foo string
	}{}
	err := response.Unmarshal(&value)

	require.NoError(t, err)
	require.Equal(t, "bar", value.Foo)
}

func TestResponseUnmarshalGzip(t *testing.T) {
	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	response := Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(gzipString(`{"foo":"bar"}`)),
			Header:     responseHeader,
		},
	}

	value := struct {
		Foo string
	}{}
	err := response.Unmarshal(&value)

	require.NoError(t, err)
	require.Equal(t, "bar", value.Foo)
}

type ReadCloseVerifier struct {
	io.Reader
	closed bool
}

func (v *ReadCloseVerifier) Close() error {
	v.closed = true
	return nil
}

func TestResponseUnmarshalClosesReader(t *testing.T) {
	stub := &ReadCloseVerifier{
		Reader: bytes.NewBufferString(`{"foo":"bar"}`),
		closed: false,
	}

	response := Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       stub,
			Header:     make(http.Header),
		},
	}

	err := response.Unmarshal(&struct{}{})

	require.NoError(t, err)
	require.True(t, stub.closed)
}

func TestResponseUnmarshalClosesInnerReader(t *testing.T) {
	stub := &ReadCloseVerifier{
		Reader: gzipString(`{"foo":"bar"}`),
		closed: false,
	}

	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	response := Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       stub,
			Header:     responseHeader,
		},
	}

	err := response.Unmarshal(&struct{}{})

	require.NoError(t, err)
	require.True(t, stub.closed)
}

func gzipString(data string) io.Reader {
	buf := new(bytes.Buffer)

	w := gzip.NewWriter(buf)
	defer w.Close()

	if _, err := io.WriteString(w, data); err != nil {
		panic(err)
	}

	return buf
}
