package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/client/responsereader"
	"github.com/SKF/go-rest-utility/problems"
)

type ProblemDecoder interface {
	DecodeProblem(context.Context, *http.Response) (problems.Problem, error)
}

type BasicProblemDecoder struct{}

func (d *BasicProblemDecoder) DecodeProblem(_ context.Context, resp *http.Response) (problems.Problem, error) {
	readBytes, err := responsereader.DecompressAndRead(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read and decompress response: %w", err)
	}

	problem := problems.BasicProblem{}
	if err := json.NewDecoder(bytes.NewReader(readBytes)).Decode(&problem); err != nil {
		return nil, fmt.Errorf("BasicProblem json decoder: %w", err)
	}

	return problem, nil
}
