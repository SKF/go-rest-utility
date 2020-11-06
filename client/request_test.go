package client

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandURL(t *testing.T) {
	baseURL := urlMustParse("https://example.com/")
	request := Get("endpoint/{id}{?limit}").
		Assign("id", 1).
		Assign("limit", "treefiddy")

	url, err := request.ExpandURL(baseURL)

	require.NoError(t, err)
	require.Equal(t, "https://example.com/endpoint/1?limit=treefiddy", url.String())
}

func TestExpandURLWithBadTemplate(t *testing.T) {
	baseURL := urlMustParse("https://example.com/")
	request := Get("endpoint/{id").Assign("id", 1)

	_, err := request.ExpandURL(baseURL)

	require.Error(t, err)
}

func TestExpandURLWithNoVariableAssignments(t *testing.T) {
	baseURL := urlMustParse("https://example.com/")
	request := Get("endpoint/{id")

	_, err := request.ExpandURL(baseURL)

	require.Error(t, err)
}

func TestExpandURLWithNoBaseURL(t *testing.T) {
	request := Get("https://example.com/endpoint/{id}").Assign("id", 1)

	url, err := request.ExpandURL(nil)

	require.NoError(t, err)
	require.Equal(t, "https://example.com/endpoint/1", url.String())
}

func urlMustParse(rawurl string) *url.URL {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return parsedURL
}
