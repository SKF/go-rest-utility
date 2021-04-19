package problems

import (
	"context"
	"net/http"
)

// Problem, basic interface for all errors supporting https://tools.ietf.org/html/rfc7807.
type Problem interface {
	error

	// ProblemType, a URI reference that identifies the problem type.
	// When dereferenced this should provide human-readable documentation for the
	// problem type. When member is not present it is assumed to be "about:blank".
	ProblemType() string

	// ProblemTitle, a short, human-readable summary of the problem type.
	// This should always be the same value for the same Type.
	ProblemTitle() string

	// ProblemStatus, the HTTP status code associated with this problem occurrence.
	// If the problem returns 0, it will be set to http.StatusInternalServerError.
	ProblemStatus() int
}

type RequestDecoratableProblem interface {
	Problem

	// DecorateWithRequest, attempt to decorate the Problem with information from the request.
	DecorateWithRequest(ctx context.Context, r *http.Request)
}

// FromError, convert an Go error into a Problem. This is a no-op if the supplied
// error already is a problem. Otherwise the returned error is an InternalProblem (HTTP 500).
func FromError(err error) Problem {
	if problem, alreadyProblem := err.(Problem); alreadyProblem {
		return problem
	}

	return Internal(err)
}
