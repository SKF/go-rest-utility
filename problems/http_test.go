package problems_test

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"

	"github.com/SKF/go-rest-utility/problems"
)

func runWriteResponse(err error, correlationID uint64) *http.Response {
	spanCtx := trace.SpanContext{}
	binary.BigEndian.PutUint64(spanCtx.TraceID[8:], correlationID)

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpanWithRemoteParent(r.Context(), "Handler", spanCtx)
		defer span.End()

		problems.WriteResponse(ctx, err, w, r)
	}

	r := httptest.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	return w.Result()
}

func TestWriteResponse_VanillaError(t *testing.T) {
	var (
		actualError                = fmt.Errorf("hello")
		actualCorrelationID uint64 = 3735928559
	)

	resp := runWriteResponse(actualError, actualCorrelationID)
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var problem problems.BasicProblem
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, "/problems/internal-server-error", problem.Type)
	require.Equal(t, "http://example.com/", problem.Instance)
	require.Equal(t, "3735928559", problem.CorrelationID)
}

type ImportantProblem struct {
	problems.BasicProblem

	Important string
}

func TestWriteResponse_DecoratableProblem(t *testing.T) {
	var (
		actualError = ImportantProblem{
			BasicProblem: problems.BasicProblem{
				Type:   "/problems/custom",
				Title:  "Custom Problem.",
				Status: http.StatusTeapot,
				Detail: "I'm very important!",
			},
			Important: "Very!",
		}
		actualCorrelationID uint64 = 3735928559
	)

	resp := runWriteResponse(actualError, actualCorrelationID)
	defer resp.Body.Close()

	require.Equal(t, http.StatusTeapot, resp.StatusCode)

	var problem ImportantProblem
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, actualError.Type, problem.Type)
	require.Equal(t, actualError.Important, problem.Important)
	require.Equal(t, "http://example.com/", problem.Instance)
	require.Equal(t, "3735928559", problem.CorrelationID)
}
