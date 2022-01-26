package auth_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/auth"
)

type SecretsManagerMock struct {
	mock.Mock
}

func (mock *SecretsManagerMock) GetSecretByID(ctx context.Context, id string) ([]byte, error) {
	args := mock.Called(ctx, id)
	return args.Get(0).([]byte), args.Error(1)
}

func TestSecretCredentialsTokenProvider_Happy(t *testing.T) {
	t.Parallel()

	var (
		user = TestUser{
			Username: "john.doe@example.com",
			Password: "very-secret",
		}
		secretID = "valid-arn"
	)

	sso := NewSSO().WithUser(user)
	defer sso.Close()

	client := new(SecretsManagerMock)
	client.On("GetSecretByID", mock.Anything, secretID).Return(
		json.Marshal(auth.SecretCredentials{
			Username: user.Username,
			Password: user.Password,
			Endpoint: sso.URL,
		}),
	).Once()

	provider := &auth.SecretCredentialsTokenProvider{
		SecretID:      secretID,
		SecretsClient: client,
	}

	token, err := provider.GetRawToken(context.Background())

	require.NoError(t, err)
	sso.RequireValidToken(t, user.Username, token)
}

func TestSecretCredentialsTokenProvider_CachesCredentials(t *testing.T) {
	t.Parallel()

	var (
		user = TestUser{
			Username: "john.doe@example.com",
			Password: "very-secret",
		}
		secretID = "valid-arn"
	)

	sso := NewSSO().WithUser(user)
	defer sso.Close()

	client := new(SecretsManagerMock)
	client.On("GetSecretByID", mock.Anything, secretID).Return(
		json.Marshal(auth.SecretCredentials{
			Username: user.Username,
			Password: user.Password,
			Endpoint: sso.URL,
		}),
	).Once()

	provider := &auth.SecretCredentialsTokenProvider{
		SecretID:      secretID,
		SecretsClient: client,
	}

	token1, err := provider.GetRawToken(context.Background())
	require.NoError(t, err)
	sso.RequireValidToken(t, user.Username, token1)

	token2, err := provider.GetRawToken(context.Background())
	require.NoError(t, err)
	sso.RequireValidToken(t, user.Username, token2)

	require.NotEqual(t, token1, token2)

	client.AssertExpectations(t)
}

func TestSecretCredentialsTokenProvider_PasswordChange(t *testing.T) {
	t.Parallel()

	var (
		user = TestUser{
			Username: "john.doe@example.com",
			Password: "very-secret",
		}
		newPassword = "new-password"
		secretID    = "valid-arn"
	)

	sso := NewSSO().WithUser(user)
	defer sso.Close()

	client := new(SecretsManagerMock)
	client.On("GetSecretByID", mock.Anything, secretID).Return(
		json.Marshal(auth.SecretCredentials{
			Username: user.Username,
			Password: user.Password,
			Endpoint: sso.URL,
		}),
	).Once()
	client.On("GetSecretByID", mock.Anything, secretID).Return(
		json.Marshal(auth.SecretCredentials{
			Username: user.Username,
			Password: newPassword,
			Endpoint: sso.URL,
		}),
	).Once()

	provider := &auth.SecretCredentialsTokenProvider{
		SecretID:      secretID,
		SecretsClient: client,
	}

	token, err := provider.GetRawToken(context.Background())
	require.NoError(t, err)
	sso.RequireValidToken(t, user.Username, token)

	user.Password = newPassword
	sso.WithUser(user)

	token, err = provider.GetRawToken(context.Background())
	require.NoError(t, err)
	sso.RequireValidToken(t, user.Username, token)
}
