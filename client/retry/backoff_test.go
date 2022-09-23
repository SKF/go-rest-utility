package retry_test

import (
	"bytes"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/retry"
)

func TestExponentialJitterBackoffProvider_FixedJitter(t *testing.T) {
	backoff := retry.ExponentialJitterBackoff{
		Base:         2 * time.Second,
		JitterSource: bytes.NewReader([]byte{0x3B, 0x9A, 0xCA, 0x00}),
	}

	actual, err := backoff.BackoffByAttempt(0)

	require.NoError(t, err)
	require.Equal(t, time.Second, actual)
}

func TestExponentialJitterBackoffProvider_Uncapped(t *testing.T) {
	provider := &retry.ExponentialJitterBackoff{
		Base: time.Second,
	}

	expectedMaximum := provider.Base

	for i := 0; i < 1_000; i++ {
		backoff, err := provider.BackoffByAttempt(i)
		require.NoError(t, err)

		require.Greater(t, backoff, time.Duration(0))
		require.Less(t, backoff, expectedMaximum)

		expectedMaximum *= 2

		if expectedMaximum < 0 {
			expectedMaximum = math.MaxInt64
		}
	}
}

func TestExponentialJitterBackoffProvider_Capped(t *testing.T) {
	provider := &retry.ExponentialJitterBackoff{
		Base: time.Second,
		Cap:  10 * time.Second,
	}

	expectedMaximum := provider.Base

	for i := 0; i < 1_000; i++ {
		backoff, err := provider.BackoffByAttempt(i)
		require.NoError(t, err)

		require.Greater(t, backoff, time.Duration(0))
		require.Less(t, backoff, expectedMaximum)

		expectedMaximum *= 2

		if expectedMaximum < provider.Cap {
			expectedMaximum = provider.Cap
		}
	}
}

func TestExponentialJitterBackoffProvider_WithMaxAttempts(t *testing.T) {
	provider := &retry.ExponentialJitterBackoff{
		MaxAttempts: 3,
	}

	for i := 0; i <= provider.MaxAttempts; i++ {
		_, err := provider.BackoffByAttempt(i)
		require.NoError(t, err)
	}

	_, err := provider.BackoffByAttempt(provider.MaxAttempts + 1)
	require.ErrorIs(t, err, retry.ErrBackoffExhausted)
}
