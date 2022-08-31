package problems

import (
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestBasicProblemDecorateWithRequest_OpenCensusCorrelationID(t *testing.T) {
	spanCtx := trace.SpanContext{}
	binary.BigEndian.PutUint64(spanCtx.TraceID[8:], uint64(3735928559))

	ctx, _ := trace.StartSpanWithRemoteParent(context.Background(), "foo", spanCtx)
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	problem := new(BasicProblem)
	problem.DecorateWithRequest(ctx, r)

	require.Equal(t, "3735928559", problem.CorrelationID)
}

func TestBasicProblemDecorateWithRequest_DatadogCorrelationID(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	span := tracer.StartSpan("foo", tracer.WithSpanID(3735928559))
	ctx := tracer.ContextWithSpan(context.Background(), span)

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	problem := new(BasicProblem)
	problem.DecorateWithRequest(ctx, r)

	require.Equal(t, "3735928559", problem.CorrelationID)
}
