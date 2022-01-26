package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/auth"
)

func TestCredentialsTokenProvider_ValidToken(t *testing.T) {
	t.Parallel()

	sso := NewSSO().WithUser(TestUser{
		Username: "john.doe@example.com",
		Password: "secret",
	})
	defer sso.Close()

	provider := auth.CredentialsTokenProvider{
		Username: "john.doe@example.com",
		Password: "secret",
		Endpoint: sso.URL,
	}

	token, err := provider.GetRawToken(context.Background())
	require.NoError(t, err)

	sso.RequireValidToken(t, "john.doe@example.com", token)
}

func TestCredentialsTokenProvider_IncorrectCredentials(t *testing.T) {
	t.Parallel()

	sso := NewSSO().WithUser(TestUser{
		Username: "john.doe@example.com",
		Password: "correct-password",
	})
	defer sso.Close()

	provider := auth.CredentialsTokenProvider{
		Username: "john.doe@example.com",
		Password: "incorrect-password",
		Endpoint: sso.URL,
	}

	_, err := provider.GetRawToken(context.Background())

	require.ErrorIs(t, err, auth.ErrIncorrectCredentials)
}

func TestCredentialsTokenProvider_Ratelimited(t *testing.T) {
	t.Parallel()

	sso := NewSSO().WithUser(TestUser{
		Username:    "john.doe@example.com",
		Password:    "very-secret",
		Ratelimited: true,
	})
	defer sso.Close()

	provider := auth.CredentialsTokenProvider{
		Username: "john.doe@example.com",
		Password: "very-secret",
		Endpoint: sso.URL,
	}

	_, err := provider.GetRawToken(context.Background())

	require.ErrorIs(t, err, auth.ErrTooManyRequests)
}

func TestCredentialsTokenProvider_Challenged(t *testing.T) {
	t.Parallel()

	sso := NewSSO().WithUser(TestUser{
		Username:   "john.doe@example.com",
		Password:   "temporary-password",
		Challenged: true,
	})
	defer sso.Close()

	provider := auth.CredentialsTokenProvider{
		Username: "john.doe@example.com",
		Password: "temporary-password",
		Endpoint: sso.URL,
	}

	_, err := provider.GetRawToken(context.Background())

	require.ErrorIs(t, err, auth.ErrChallenged)
}

func TestCredentialsTokenProvider_Inactive(t *testing.T) {
	t.Parallel()

	sso := NewSSO().WithUser(TestUser{
		Username: "john.doe@example.com",
		Password: "very-secret",
		Inactive: true,
	})
	defer sso.Close()

	provider := auth.CredentialsTokenProvider{
		Username: "john.doe@example.com",
		Password: "very-secret",
		Endpoint: sso.URL,
	}

	_, err := provider.GetRawToken(context.Background())

	require.ErrorIs(t, err, auth.ErrInactivated)
}

func TestCredentialsTokenProvider_UnknownTokenType(t *testing.T) {
	t.Parallel()

	sso := NewSSO().WithUser(TestUser{
		Username: "john.doe@example.com",
		Password: "very-secret",
	})
	defer sso.Close()

	provider := auth.CredentialsTokenProvider{
		Username:  "john.doe@example.com",
		Password:  "very-secret",
		Endpoint:  sso.URL,
		TokenType: "unknown",
	}

	_, err := provider.GetRawToken(context.Background())

	require.ErrorIs(t, err, auth.ErrUnknownTokenType)
}
