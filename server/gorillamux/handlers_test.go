package gorillamux

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SKF/go-rest-utility/server/gorillamux/problems"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MethodNotFoundHandler(t *testing.T) {
	expectedRes := problems.MethodNotAllowed("GET", "PUT", "POST")

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
	assert.Equal(t, expectedRes, actual)

}

func Test_NotFoundHandler(t *testing.T) {
	expectedRes := problems.NotFound()

	ts := httptest.NewServer(NotFoundHandler())
	defer ts.Close()
	res, err := http.Get(ts.URL)

	require.NoError(t, err)

	reader := res.Body
	defer reader.Close()

	actual := problems.NotFoundProblem{}

	require.NoError(t, json.NewDecoder(reader).Decode(&actual))
	assert.Equal(t, expectedRes, actual)
}
