package auth_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/auth"
)

type TokenProviderMock struct {
	mock.Mock
}

func (mock *TokenProviderMock) GetRawToken(ctx context.Context) (auth.RawToken, error) {
	args := mock.Called(ctx)
	return args.Get(0).(auth.RawToken), args.Error(1)
}

func testProviderInParallel(t *testing.T, provider auth.TokenProvider, n int, condition func(auth.RawToken, error)) {
	t.Helper()

	// How many of the "n" parallell runs that should be invoked with a small delay (1s) before
	// calling GetRawToken. Set to the last 25% of invocations
	staggerLimit := 3 * n / 4

	wg := new(sync.WaitGroup)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(p auth.TokenProvider, i int) {
			defer wg.Done()

			// Stagger invocations to inc
			if i > staggerLimit {
				time.Sleep(1 * time.Second)
			}

			condition(p.GetRawToken(context.Background()))
		}(provider, i)
	}

	wg.Wait()
}

func TestCachedTokenProvider_OnlyOneSupplierCall(t *testing.T) {
	t.Parallel()

	expectedToken := TestAccessToken{
		Email:    "john.doe@example.com",
		Lifetime: 1 * time.Hour,
	}.Build(t)

	provider := new(TokenProviderMock)
	provider.On("GetRawToken", mock.Anything).After(2*time.Second).Return(expectedToken, nil).Once()

	cached := auth.NewCachedTokenProvider(provider)

	testProviderInParallel(t, cached, 2000, func(token auth.RawToken, err error) {
		require.NoError(t, err)
		require.Equal(t, expectedToken, token)
	})

	provider.AssertExpectations(t)
}

func TestCachedTokenProvider_AlwaysFailingProvider(t *testing.T) {
	t.Parallel()

	expectedError := fmt.Errorf("FOOBAR")

	provider := new(TokenProviderMock)
	provider.On("GetRawToken", mock.Anything).After(100*time.Millisecond).Return(auth.RawToken(""), expectedError)

	cached := auth.NewCachedTokenProvider(provider)

	testProviderInParallel(t, cached, 10, func(token auth.RawToken, err error) {
		require.ErrorIs(t, err, expectedError)
	})

	// We expect the calls to be between 1 and 10, depending on the scheduler.
	require.LessOrEqual(t, len(provider.Calls), 10)
	require.Greater(t, len(provider.Calls), 0)

	provider.AssertExpectations(t)
}

func TestCachedTokenProvider_DoubleCaching(t *testing.T) {
	t.Parallel()

	provider := new(TokenProviderMock)

	wrappedOnce := auth.NewCachedTokenProvider(provider)
	wrappedTwice := auth.NewCachedTokenProvider(wrappedOnce)

	require.Same(t, wrappedOnce, wrappedTwice)
}

func TestCachedTokenProvider_InvalidToken(t *testing.T) {
	t.Parallel()

	provider := new(TokenProviderMock)
	provider.On("GetRawToken", mock.Anything).Return(auth.RawToken("bad-token"), nil)

	cached := auth.NewCachedTokenProvider(provider)

	_, err := cached.GetRawToken(context.Background())

	require.Error(t, err)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestCachedTokenProvider_RefreshAfterTTL(t *testing.T) {
	t.Parallel()

	lifetime := 10 * time.Second
	gracePeriod := 1 * time.Second

	fc := &FakeClock{
		Now: time.Now(),
	}

	expectedToken1 := TestAccessToken{
		Email:     "john.doe@example.com",
		Lifetime:  lifetime,
		IssueTime: fc.Now,
	}.Build(t)
	expectedToken2 := TestAccessToken{
		Email:     "john.doe@example.com",
		Lifetime:  lifetime,
		IssueTime: fc.Now.Add(lifetime).Add(-gracePeriod),
	}.Build(t)

	ctx := context.Background()

	provider := new(TokenProviderMock)
	provider.On("GetRawToken", ctx).Return(expectedToken1, nil).Once()
	provider.On("GetRawToken", ctx).Return(expectedToken2, nil).Once()

	cached := auth.NewCachedTokenProvider(provider).
		WithGracePeriod(gracePeriod).
		WithClock(fc.Get)

	actualToken1, err := cached.GetRawToken(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedToken1, actualToken1)

	fc.Now = fc.Now.Add(time.Second)

	actualToken2, err := cached.GetRawToken(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedToken1, actualToken2)

	// Set the current time to within the lifetime of the last token, but within the grace period
	fc.Now = fc.Now.Add(lifetime).Add(-time.Second).Add(-gracePeriod)

	actualToken3, err := cached.GetRawToken(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedToken2, actualToken3)

	provider.AssertExpectations(t)
}
