package client

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ReadCloseVerifier struct {
	io.Reader
	closed bool
}

func (v *ReadCloseVerifier) Close() error {
	v.closed = true
	return nil
}

func TestResponseUnmarshalSimple(t *testing.T) {
	response := Response{
		Response: http.Response{
			StatusCode: http.StatusOK,
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

func TestDecompressResponse(t *testing.T) {
	response := http.Response{ //nolint:bodyclose
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
		Header:     make(http.Header),
	}

	body, header, err := DecompressResponse(response)

	require.NoError(t, err)
	readBytes, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	assert.Equal(t, `{"foo":"bar"}`, string(readBytes))
	assert.Equal(t, "", header.Get(headers.ContentEncoding))
}

func TestDecompressResponseGzip(t *testing.T) {
	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	response := http.Response{ //nolint:bodyclose
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(gzipString(`{"foo":"bar"}`)),
		Header:     responseHeader,
	}
	body, header, err := DecompressResponse(response)

	require.NoError(t, err)
	readBytes, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	assert.Equal(t, `{"foo":"bar"}`, string(readBytes))
	assert.Equal(t, "", header.Get(headers.ContentEncoding))
}

func TestDecompressResponseGzipButContentLengthZero(t *testing.T) {
	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")
	responseHeader.Set(headers.ContentLength, "0")

	response := http.Response{ //nolint:bodyclose
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(``)),
		Header:     responseHeader,
	}
	body, header, err := DecompressResponse(response)

	require.NoError(t, err)
	readBytes, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	assert.Equal(t, ``, string(readBytes))
	assert.Equal(t, "gzip", header.Get(headers.ContentEncoding))
}

func TestDecompressResponseGzipInnerBodyIsClosed(t *testing.T) {
	verifier := ReadCloseVerifier{
		Reader: gzipString(`{"foo":"bar"}`),
		closed: false,
	}

	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	response := http.Response{
		StatusCode: http.StatusOK,
		Body:       &verifier,
		Header:     responseHeader,
	}

	_, _, err := DecompressResponse(response)
	require.NoError(t, err)

	_ = response.Body.Close()

	require.True(t, verifier.closed)
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
