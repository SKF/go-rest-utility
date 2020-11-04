package auth

import (
	"context"
	"fmt"

	"github.com/SKF/go-utility/v2/auth/secretsmanagerauth"
)

type IdentityToken string

type TokenProvider interface {
	GetIdentityToken(ctx context.Context) (IdentityToken, error)
}

func (token IdentityToken) GetIdentityToken(ctx context.Context) (IdentityToken, error) {
	return token, nil
}

func (token IdentityToken) String() string {
	return string(token)
}

type SecretsManagerTokenProvider struct {
	configured bool
	Config     secretsmanagerauth.Config
}

func (provider *SecretsManagerTokenProvider) GetIdentityToken(ctx context.Context) (IdentityToken, error) {
	if !provider.configured {
		secretsmanagerauth.Configure(provider.Config)
		provider.configured = true
	}

	if err := secretsmanagerauth.SignIn(ctx); err != nil {
		return "", fmt.Errorf("unable to sign-in as service user '%s': %w", provider.Config.SecretKey, err)
	}

	return IdentityToken(secretsmanagerauth.GetTokens().IdentityToken), nil
}
