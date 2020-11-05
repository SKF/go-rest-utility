package client

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
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

	if assert.NoError(t, err) {
		assert.Equal(t, "bar", value.Foo)
	}
}

func TestResponseUnmarshalGzip(t *testing.T) {
	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	response := Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       gzipString(`{"foo":"bar"}`),
			Header:     responseHeader,
		},
	}

	value := struct {
		Foo string
	}{}
	err := response.Unmarshal(&value)

	if assert.NoError(t, err) {
		assert.Equal(t, "bar", value.Foo)
	}
}

func gzipString(data string) io.ReadCloser {
	buf := new(bytes.Buffer)

	w := gzip.NewWriter(buf)
	defer w.Close()

	if _, err := io.WriteString(w, data); err != nil {
		panic(err)
	}

	return ioutil.NopCloser(buf)
}
