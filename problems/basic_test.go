package problems

import (
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	opencensus "go.opencensus.io/trace"
	datadog_mock "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	datadog "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestBasicProblemDecorateWithRequest_OpenCensusCorrelationID(t *testing.T) {
	spanCtx := opencensus.SpanContext{}
	binary.BigEndian.PutUint64(spanCtx.TraceID[8:], uint64(3735928559))

	ctx, _ := opencensus.StartSpanWithRemoteParent(context.Background(), "foo", spanCtx)
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	problem := new(BasicProblem)
	problem.DecorateWithRequest(ctx, r)

	require.Equal(t, "3735928559", problem.CorrelationID)
}

func TestBasicProblemDecorateWithRequest_DatadogCorrelationID(t *testing.T) {
	mt := datadog_mock.Start()
	defer mt.Stop()

	span := datadog.StartSpan("foo", datadog.WithSpanID(3735928559))
	ctx := datadog.ContextWithSpan(context.Background(), span)

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	problem := new(BasicProblem)
	problem.DecorateWithRequest(ctx, r)

	require.Equal(t, "3735928559", problem.CorrelationID)
}
