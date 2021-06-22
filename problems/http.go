package problems

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

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
	problem = tryDecorateWithRequest(ctx, problem, r)

	statusCode := problem.ProblemStatus()

	l := log.
		WithTracing(ctx).
		WithUserID(ctx).
		WithError(err).
		WithField("code", statusCode)

	// Log as an Error if statusCode is 5XX, Warn if context was cancelled, otherwise Info.
	switch {
	case statusCode/100 == http.StatusInternalServerError/100:
		l.Error(problem.ProblemTitle())
	case ctx.Err() == context.Canceled:
		l.Warn(problem.ProblemTitle())
	default:
		l.Info(problem.ProblemTitle())
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(problem); encodeErr != nil {
		l.WithError(encodeErr).Error("Unable to write problem output")
	}
}

// tryDecorateWithRequest attempts to call DecorateWithRequest if the supplied problem implements it
//
// The returned value is the decorated problem or the input problem if not supported
func tryDecorateWithRequest(ctx context.Context, problem Problem, r *http.Request) Problem {
	problemValue := reflect.Indirect(reflect.ValueOf(problem))

	problemCopy := reflect.New(problemValue.Type())
	problemCopy.Elem().Set(problemValue)

	if decoratableProblem, ok := problemCopy.Interface().(RequestDecoratableProblem); ok {
		decoratableProblem.DecorateWithRequest(ctx, r)

		return decoratableProblem
	}

	return problem
}
