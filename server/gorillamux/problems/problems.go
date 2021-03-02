package problems

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SKF/go-rest-utility/problems"
)

type NotFoundProblem struct {
	problems.BasicProblem
}

func NotFound() NotFoundProblem {
	return NotFoundProblem{
		BasicProblem: problems.BasicProblem{
			Type:   "/problems/route-not-found",
			Title:  "The requested endpoint could not be found.",
			Status: http.StatusNotFound,
			Detail: "Ensure that the URI is a valid endpoint for the service.",
		},
	}
}

type MethodNotAllowedProblem struct {
	problems.BasicProblem
	Method  string   `json:"requested_method,omitempty"`
	Allowed []string `json:"allowed_methods,omitempty"`
}

func MethodNotAllowed(requested string, allowed ...string) MethodNotAllowedProblem {
	return MethodNotAllowedProblem{
		BasicProblem: problems.BasicProblem{
			Type:   "/problems/request-method-not-allowed",
			Title:  "The requested method is not allowed.",
			Status: http.StatusMethodNotAllowed,
			Detail: fmt.Sprintf(
				"The requested resource does not support method '%s', it does only support one of '%s'.",
				requested,
				strings.Join(allowed, ", "),
			),
		},
		Method:  requested,
		Allowed: allowed,
	}
}
