package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/problems"
)

type ProblemDecoder interface {
	DecodeProblem(context.Context, *http.Response) (problems.Problem, error)
}

type BasicProblemDecoder struct{}

func (d *BasicProblemDecoder) DecodeProblem(ctx context.Context, resp *http.Response) (problems.Problem, error) {
	defer resp.Body.Close()

	problem := problems.BasicProblem{}
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		return nil, fmt.Errorf("BasicProblem json decoder: %w", err)
	}

	return problem, nil
}
