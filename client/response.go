package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/client/responsereader"
)

type Response struct {
	http.Response
}

func (r *Response) Unmarshal(v interface{}) error {
	readBytes, err := responsereader.DecompressAndRead(&r.Response)
	if err != nil {
		return fmt.Errorf("failed to read and decompress response: %w", err)
	}

	if err = json.NewDecoder(bytes.NewReader(readBytes)).Decode(&v); err != nil {
		return fmt.Errorf("failed to json decode read bytes: %w", err)
	}

	return nil
}
