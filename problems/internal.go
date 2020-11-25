package problems

import (
	"errors"
	"fmt"
	"net/http"
)

type InternalProblem struct {
	BasicProblem
	cause error `json:"-"`
}

func Internal(cause error) InternalProblem {
	if cause == nil {
		cause = errors.New("nil internal server error")
	}

	return InternalProblem{
		BasicProblem: BasicProblem{
			Type:   "/problems/internal-server-error",
			Title:  "Internal Server Error",
			Status: http.StatusInternalServerError,
			Detail: "An unexpected error occurred, please contact support or try again.",
		},
		cause: cause,
	}
}

func (problem InternalProblem) Error() string {
	return problem.cause.Error()
}

func (problem InternalProblem) Unwrap() error {
	return problem.cause
}

func (problem InternalProblem) Format(f fmt.State, c rune) {
	if errFmt, ok := problem.cause.(fmt.Formatter); ok {
		errFmt.Format(f, c)
	} else {
		fmt.Fprintf(f, "%s", problem.cause.Error())
	}
}
