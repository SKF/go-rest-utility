package problems

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromError_vanillaError_InternalProblem(t *testing.T) {
	// Given
	err := fmt.Errorf("vanilla")

	// When
	result := FromError(err)

	_, ok := result.(InternalProblem)
	require.True(t, ok)
}

func TestFromError_ValidationProblem_ValidationProblem(t *testing.T) {
	// Given
	err := ValidationProblem{}

	// When
	result := FromError(err)

	_, ok := result.(ValidationProblem)
	require.True(t, ok)
}

func TestFromError_WrappedValidationProblem_ValidationProblem(t *testing.T) {
	// Given
	err := fmt.Errorf("wrapped this problem: %w", ValidationProblem{})

	// When
	result := FromError(err)

	_, ok := result.(ValidationProblem)
	require.True(t, ok)
}
