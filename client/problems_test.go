package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/problems"
)

// The default client should not consider Problems for backward compatibility reasons
// should probably be included in the default client in the next major version (v1).
func TestClientGetWithoutProblemDecoder(t *testing.T) {
	_, err := setupProblemServer(errors.New("internal error"))

	require.Error(t, err)
	require.ErrorIs(t, err, client.ErrInternalServerError)
}

func TestClientGetWithBasicProblemDecoder(t *testing.T) {
	expectedProblem := problems.BasicProblem{
		Type:   "/my-basic-problem",
		Title:  "A Basic Problem",
		Status: http.StatusTeapot,
		Detail: "What did you expect?",
	}

	_, err := setupProblemServer(expectedProblem, client.WithProblemDecoder(new(client.BasicProblemDecoder)))
	require.Error(t, err)
	require.Implements(t, (*problems.Problem)(nil), err)

	actualProblem := err.(problems.Problem)
	require.Equal(t, expectedProblem.Type, actualProblem.ProblemType())
	require.Equal(t, expectedProblem.Status, actualProblem.ProblemStatus())
	require.Equal(t, expectedProblem.Title, actualProblem.ProblemTitle())
}

func TestClientGetWithCustomProblemDecoder(t *testing.T) {
	type CustomProblem struct {
		problems.BasicProblem
		Counter int
	}

	expectedProblem := CustomProblem{
		BasicProblem: problems.BasicProblem{
			Type:   "/my-custom-problem",
			Title:  "An Custom Problem",
			Status: http.StatusTeapot,
			Detail: "What did you expect?",
		},
		Counter: 10,
	}

	var decoder ProblemDecoderFn = func(ctx context.Context, resp *http.Response) (problems.Problem, error) {
		defer resp.Body.Close()

		problem := CustomProblem{}
		if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
			return nil, err
		}

		return problem, nil
	}

	_, err := setupProblemServer(expectedProblem, client.WithProblemDecoder(decoder))
	require.Error(t, err)
	require.IsType(t, expectedProblem, err)

	actualProblem := err.(CustomProblem)
	require.Equal(t, expectedProblem.Type, actualProblem.Type)
	require.Equal(t, expectedProblem.Title, actualProblem.Title)
	require.Equal(t, expectedProblem.Status, actualProblem.Status)
	require.Equal(t, expectedProblem.Detail, actualProblem.Detail)
	require.Equal(t, expectedProblem.Counter, actualProblem.Counter)
}

func setupProblemServer(problem error, opts ...client.Option) (*client.Response, error) { //nolint:unparam
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		problems.WriteResponse(r.Context(), problem, w, r)
	}))
	defer srv.Close()

	request := client.Get("endpoint")

	client := client.NewClient(
		append([]client.Option{
			client.WithBaseURL(srv.URL),
		}, opts...)...,
	)

	return client.Do(context.Background(), request)
}

type ProblemDecoderFn func(context.Context, *http.Response) (problems.Problem, error)

func (fn ProblemDecoderFn) DecodeProblem(ctx context.Context, resp *http.Response) (problems.Problem, error) {
	return fn(ctx, resp)
}
