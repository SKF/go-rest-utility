package problems

import (
	"context"
	"encoding/binary"
	"fmt"
	"net/http"
	"strconv"

	"go.opencensus.io/trace"
)

// BasicProblem, common fields for most Problems. Useful for embedding into
// for custom problem types.
type BasicProblem struct {
	// Type, a URI reference that identifies the problem type.
	// When dereferenced this should provide human-readable documentation for the
	// problem type. When member is not present it is assumed to be "about:blank".
	Type string `json:"type" format:"url"`

	// Title, a short, human-readable summary of the problem type.
	// This should always be the same value for the same Type.
	Title string `json:"title"`

	// Status, the HTTP status code associated with this problem occurrence.
	Status int `json:"status,omitempty"`

	// Detail, a human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// Instance, a URI reference that identifies the specific resource on which the problem occurred.
	Instance string `json:"instance,omitempty" format:"url"`

	// CorrelationID, an unique identifier for tracing this issue in server logs.
	CorrelationID string `json:"correlation_id,omitempty"`
}

// Ensure that BasicProblem implements all methods required by Problem.
var _ Problem = BasicProblem{}

func (problem BasicProblem) ProblemType() string {
	if problem.Type == "" {
		return "about:blank"
	}

	return problem.Type
}

func (problem BasicProblem) ProblemTitle() string {
	return problem.Title
}

func (problem BasicProblem) ProblemStatus() int {
	if problem.Status <= 0 {
		return http.StatusInternalServerError
	}

	return problem.Status
}

func (problem BasicProblem) Error() string {
	if problem.Detail == "" {
		return problem.ProblemTitle()
	}

	return fmt.Sprintf("%s: %s", problem.ProblemTitle(), problem.Detail)
}

func (problem *BasicProblem) DecorateWithRequest(ctx context.Context, r *http.Request) {
	// Extract the DataDog TraceID from the request context. Same logic can be found in
	// the go-utility/log package in the WithTracing function.
	if span := trace.FromContext(ctx); span != nil {
		traceID := span.SpanContext().TraceID
		problem.CorrelationID = strconv.FormatUint(binary.BigEndian.Uint64(traceID[8:]), 10)
	}

	problem.Instance = r.URL.String()
}
