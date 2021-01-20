package problems_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/problems"
)

func TestAppendWithSimplePrefix(t *testing.T) {
	problem := problems.Validation()

	problem.AppendWithPrefix("a", problems.ValidationReason{
		Name:   "b",
		Reason: "c",
	})

	require.Len(t, problem.Reasons, 1)
	require.Equal(t, "a.b", problem.Reasons[0].Name)
	require.Equal(t, "c", problem.Reasons[0].Reason)
}

func TestAppendWithEmptyPrefix(t *testing.T) {
	problem := problems.Validation()

	problem.AppendWithPrefix("", problems.ValidationReason{
		Name:   "b",
		Reason: "c",
	})

	require.Len(t, problem.Reasons, 1)
	require.Equal(t, "b", problem.Reasons[0].Name)
	require.Equal(t, "c", problem.Reasons[0].Reason)
}

func TestAppendWithPrefixAndNoName(t *testing.T) {
	problem := problems.Validation()

	problem.AppendWithPrefix("a", problems.ValidationReason{
		Reason: "c",
	})

	require.Len(t, problem.Reasons, 1)
	require.Equal(t, "a", problem.Reasons[0].Name)
	require.Equal(t, "c", problem.Reasons[0].Reason)
}
