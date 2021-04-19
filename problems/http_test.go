package problems_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"

	"github.com/SKF/go-rest-utility/problems"
)

func runWriteResponse(err error) *http.Response {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "Handler")
		defer span.End()

		problems.WriteResponse(ctx, err, w, r)
	}

	r := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	return w.Result()
}

func TestWriteResponse_VanillaError(t *testing.T) {
	actualError := fmt.Errorf("hello")

	resp := runWriteResponse(actualError)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var problem problems.BasicProblem
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, "/problems/internal-server-error", problem.Type)
	require.Equal(t, "http://example.com/", problem.Instance)
	require.NotEmpty(t, problem.CorrelationID)
}

type ImportantProblem struct {
	problems.BasicProblem

	Important string
}

func TestWriteResponse_DecoratableProblem(t *testing.T) {
	actualError := ImportantProblem{
		BasicProblem: problems.BasicProblem{
			Type:   "/problems/custom",
			Title:  "Custom Problem.",
			Status: http.StatusTeapot,
			Detail: "I'm very important!",
		},
		Important: "Very!",
	}

	resp := runWriteResponse(actualError)
	defer resp.Body.Close()

	require.Equal(t, http.StatusTeapot, resp.StatusCode)

	var problem ImportantProblem
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, actualError.Type, problem.Type)
	require.Equal(t, actualError.Important, problem.Important)
	require.Equal(t, "http://example.com/", problem.Instance)
	require.NotEmpty(t, problem.CorrelationID)
}
