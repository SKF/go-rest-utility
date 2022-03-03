package responsereader_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/responsereader"
)

type ReadCloseVerifier struct {
	io.Reader
	closed bool
}

func (v *ReadCloseVerifier) Close() error {
	v.closed = true
	return nil
}

func TestDecompressAndRead(t *testing.T) {
	readBytes, err := responsereader.DecompressAndRead(&http.Response{ //nolint:bodyclose
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
		Header:     make(http.Header),
	})

	require.NoError(t, err)
	require.Equal(t, `{"foo":"bar"}`, string(readBytes))
}

func TestHandleResponseGzip(t *testing.T) {
	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	readBytes, err := responsereader.DecompressAndRead(&http.Response{ //nolint:bodyclose
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(gzipString(`{"foo":"bar"}`)),
		Header:     responseHeader,
	})

	require.NoError(t, err)
	require.Equal(t, `{"foo":"bar"}`, string(readBytes))
}

func TestHandleResponseGzipInnerBodyIsClosed(t *testing.T) {
	verifier := ReadCloseVerifier{
		Reader: gzipString(`{"foo":"bar"}`),
		closed: false,
	}

	responseHeader := make(http.Header)
	responseHeader.Set(headers.ContentEncoding, "gzip")

	_, err := responsereader.DecompressAndRead(&http.Response{ //nolint:bodyclose
		StatusCode: http.StatusOK,
		Body:       &verifier,
		Header:     responseHeader,
	})

	require.NoError(t, err)
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
