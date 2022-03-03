package responsereader

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-http-utils/headers"
)

func DecompressAndRead(response *http.Response) ([]byte, error) {
	defer response.Body.Close()
	reader := response.Body

	switch response.Header.Get(headers.ContentEncoding) {
	case "gzip":
		var err error
		if reader, err = gzip.NewReader(response.Body); err != nil {
			return nil, fmt.Errorf("failed to interpret response body as a gzip reader: %w", err)
		}
		defer reader.Close()
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read all content in response body: %w", err)
	}

	return bodyBytes, nil
}
