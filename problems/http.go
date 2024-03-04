package problems

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/SKF/go-utility/v2/log"
)

const ContentType = "application/problem+json"

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

	var allProblemFields map[string]interface{}
	marshaledProblem, _ := json.Marshal(problem)
	json.Unmarshal(marshaledProblem, &allProblemFields)

	l := log.
		WithTracing(ctx).
		WithClientID(ctx).
		WithUserID(ctx).
		WithError(err).
		WithField("problem", allProblemFields).
		WithField("code", statusCode)

	// Log as an Error if statusCode is 5XX, otherwise as Info.
	if statusCode/100 == http.StatusInternalServerError/100 {
		l.Error(problem.ProblemTitle())
	} else {
		l.Info(problem.ProblemTitle())
	}

	w.Header().Set("Content-Type", ContentType)
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
