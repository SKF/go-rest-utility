package problems

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
	}{
		{
				name: "Problem is Problem",
			err:            BasicProblem{Status: http.StatusBadRequest},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Wrapped problem is Problem",
			err:            fmt.Errorf("some text: %w", BasicProblem{Status: http.StatusBadRequest}),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Non problem returns internal server error problem",
			err:            fmt.Errorf("some text"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromError(tt.err)
			require.Equal(t, tt.wantStatusCode, got.ProblemStatus())
		})
	}
}
