package auth

import (
	"context"
	"fmt"

	"github.com/SKF/go-utility/v2/auth/secretsmanagerauth"
)

type RawToken string

type TokenProvider interface {
	GetRawToken(ctx context.Context) (RawToken, error)
}

func (token RawToken) GetRawToken(ctx context.Context) (RawToken, error) {
	return token, nil
}

func (token RawToken) String() string {
	return string(token)
}

type SecretsManagerTokenProvider struct {
	configured bool
	Config     secretsmanagerauth.Config
}

func (provider *SecretsManagerTokenProvider) GetIdentityToken(ctx context.Context) (RawToken, error) {
	if !provider.configured {
		secretsmanagerauth.Configure(provider.Config)
		provider.configured = true
	}

	if err := secretsmanagerauth.SignIn(ctx); err != nil {
		return "", fmt.Errorf("unable to sign-in as service user '%s': %w", provider.Config.SecretKey, err)
	}

	return RawToken(secretsmanagerauth.GetTokens().IdentityToken), nil
}
