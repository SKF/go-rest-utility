package client

import (
	netUrl "net/url"
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

func TestExpandURLWithEmptySlice(t *testing.T) {
	baseURL := urlMustParse("https://example.net/")
	request := Get("/person/albums{?field*}").
		Assign("field", []string{})

	url, err := request.ExpandURL(baseURL)

	require.NoError(t, err)
	require.Equal(t, "https://example.net/person/albums", url.String())
}

func TestExpandURLWithStringSlice(t *testing.T) {
	baseURL := urlMustParse("https://example.net/")
	request := Get("/person/albums{?field*}").
		Assign("field", []string{"id", "name", "picture"})

	url, err := request.ExpandURL(baseURL)

	require.NoError(t, err)
	require.Equal(t, "https://example.net/person/albums?field=id&field=name&field=picture", url.String())
}

func TestExpandURLWithStringMap(t *testing.T) {
	baseURL := urlMustParse("https://example.net/")
	request := Get("/person/albums{?keys}").
		Assign("keys", map[string]string{
			"semi":  ";",
			"dot":   ".",
			"comma": ",",
		})

	url, err := request.ExpandURL(baseURL)

	require.NoError(t, err)

	// Iteration order of Go maps is random so the expected result is all permutations
	require.Contains(t, []string{
		"https://example.net/person/albums?keys=semi,%3B,dot,.,comma,%2C",
		"https://example.net/person/albums?keys=semi,%3B,comma,%2C,dot,.",
		"https://example.net/person/albums?keys=dot,.,semi,%3B,comma,%2C",
		"https://example.net/person/albums?keys=dot,.,comma,%2C,semi,%3B",
		"https://example.net/person/albums?keys=comma,%2C,semi,%3B,dot,.",
		"https://example.net/person/albums?keys=comma,%2C,dot,.,semi,%3B",
	}, url.String())
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

func urlMustParse(rawurl string) *netUrl.URL {
	parsedURL, err := netUrl.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	return parsedURL
}
