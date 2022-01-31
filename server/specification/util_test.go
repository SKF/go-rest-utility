package specification_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/server/specification"
)

func TestBaseURLFromRequest(t *testing.T) {
	expected := "http://acme.inc"

	r := httptest.NewRequest(http.MethodGet, "/foo.txt", nil)
	r.Host = "acme.inc"

	actual := specification.BaseURLFromRequest(r)
	require.Equal(t, expected, actual.String())
}

func TestBaseURLFromRequest_WithForwardedProto(t *testing.T) {
	expected := "https://acme.inc:8080"

	r := httptest.NewRequest(http.MethodGet, "/foo.txt", nil)
	r.Host = "acme.inc:8080"
	r.Header.Add("X-Forwarded-Proto", "https")

	actual := specification.BaseURLFromRequest(r)
	require.Equal(t, expected, actual.String())
}

func TestContentTypeFromFilename_YAML(t *testing.T) {
	expected := "application/x-yaml"

	actual := specification.ContentTypeFromFilename("test.yaml")

	require.Equal(t, expected, actual)
}
