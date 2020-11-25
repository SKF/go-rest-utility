package problems

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SKF/go-utility/v2/log"
)

// Generic, returns a generic HTTP-based Problem from a HTTP status code.
func Generic(status int) Problem {
	return BasicProblem{
		Title:  http.StatusText(status),
		Status: status,
	}
}

// WriteResponse, converts the error into a Problem and writes the contents into w.
// The Problem will be decorated with request information and logged if possible.
func WriteResponse(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) {
	problem := FromError(err)

	if decoratableProblem, ok := problem.(RequestDecoratableProblem); ok {
		decoratableProblem.DecorateWithRequest(ctx, r)
	}

	statusCode := problem.ProblemStatus()

	l := log.
		WithTracing(ctx).
		WithUserID(ctx).
		WithError(err).
		WithField("code", statusCode)

	// Log as an Error if statusCode is 5XX, otherwise as Info.
	if statusCode/100 == http.StatusInternalServerError/100 {
		l.Error(problem.ProblemTitle())
	} else {
		l.Info(problem.ProblemTitle())
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(problem); encodeErr != nil {
		l.WithError(encodeErr).Error("Unable to write problem output")
	}
}
