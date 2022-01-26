package cachedtokenprovider_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	auth_model "github.com/SKF/go-utility/v2/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/auth/cachedtokenprovider"
)

const (
	expectedToken = "my-token"
)

func Test_GetRawToken_HappyCase(t *testing.T) {
	// Given
	ctx := context.Background()
	provider := cachedtokenprovider.NewWithCustomAuth(cachedtokenprovider.Config{}, &AuthStub{})

	// When
	token, err := provider.GetRawToken(ctx)

	// Then
	require.NoError(t, err)
	assert.Equal(t, "my-token.0", string(token))
}

func Test_GetRawToken_NoNewTokenWithinTTL(t *testing.T) {
	// Given
	ctx := context.Background()
	provider := cachedtokenprovider.NewWithCustomAuth(cachedtokenprovider.Config{
		TokenTimeToLive: time.Minute,
	}, &AuthStub{})

	for i := 0; i < 10; i++ {
		// When
		token, err := provider.GetRawToken(ctx)

		require.NoError(t, err)
		assert.Equal(t, "my-token.0", string(token))
	}
}

func Test_GetRawToken_NewTokenAfterTTL(t *testing.T) {
	// Given
	ctx := context.Background()
	provider := cachedtokenprovider.NewWithCustomAuth(cachedtokenprovider.Config{
		TokenTimeToLive: time.Millisecond,
	}, &AuthStub{})

	tokenBeforeTTL, err := provider.GetRawToken(ctx)
	require.NoError(t, err)
	assert.Equal(t, "my-token.0", string(tokenBeforeTTL))

	time.Sleep(2 * time.Millisecond)

	// When
	token, err := provider.GetRawToken(ctx)

	// Then
	require.NoError(t, err)
	assert.Equal(t, "my-token.1", string(token))
}

type AuthStub struct {
	i int
}

func (stub *AuthStub) GetTokens(context.Context) (auth_model.Tokens, error) {
	tokens := auth_model.Tokens{
		IdentityToken: fmt.Sprintf("%s.%d", expectedToken, stub.i),
	}

	stub.i++

	return tokens, nil
}
