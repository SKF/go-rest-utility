package problems

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SKF/go-rest-utility/problems"
)

func Test_MethodNotFoundHandler(t *testing.T) {
	expectedProblem := MethodNotAllowedProblem{
		BasicProblem: problems.BasicProblem{
			Type:   "/problems/request-method-not-allowed",
			Title:  "The requested method is not allowed.",
			Status: http.StatusMethodNotAllowed,
			Detail: fmt.Sprintf(
				"The requested resource does not support method '%s', it does only support one of '%s'.",
				"GET",
				strings.Join([]string{"PUT", "PATCH"}, ", "),
			),
		},
		Method:  "GET",
		Allowed: []string{"PUT", "PATCH"},
	}

	assert.Equal(t, expectedProblem, MethodNotAllowed("GET", "PUT", "PATCH"))
}

func Test_NotFoundHandler(t *testing.T) {
	expectedProblem := NotFoundProblem{
		BasicProblem: problems.BasicProblem{
			Type:   "/problems/route-not-found",
			Title:  "The requested endpoint could not be found.",
			Status: http.StatusNotFound,
			Detail: "Ensure that the URI is a valid endpoint for the service.",
		},
	}

	assert.Equal(t, expectedProblem, NotFound())
}
