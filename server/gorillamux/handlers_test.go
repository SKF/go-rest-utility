package gorillamux

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/server/gorillamux/problems"
)

func Test_MethodNotFoundHandler(t *testing.T) {
	expected := problems.MethodNotAllowed("GET", "PUT", "POST")

	router := mux.NewRouter()
	router.Name("testRouter").Methods("PUT", "POST")

	ts := httptest.NewServer(MethodNotFoundHandler(router))
	defer ts.Close()
	res, err := http.Get(ts.URL)

	require.NoError(t, err)

	reader := res.Body
	defer reader.Close()

	actual := problems.MethodNotAllowedProblem{}

	require.NoError(t, json.NewDecoder(reader).Decode(&actual))

	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Method, actual.Method)
	assert.Equal(t, expected.Allowed, actual.Allowed)
}

func Test_NotFoundHandler(t *testing.T) {
	expected := problems.NotFound()

	ts := httptest.NewServer(NotFoundHandler())
	defer ts.Close()
	res, err := http.Get(ts.URL)

	require.NoError(t, err)

	reader := res.Body
	defer reader.Close()

	actual := problems.NotFoundProblem{}

	require.NoError(t, json.NewDecoder(reader).Decode(&actual))

	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Status, actual.Status)
}
