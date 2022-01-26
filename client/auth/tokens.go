package auth

import (
	"context"
	"fmt"

	"github.com/SKF/go-utility/v2/auth/secretsmanagerauth"
)

type TokenProvider interface {
	GetRawToken(ctx context.Context) (RawToken, error)
}

// Deprecated: Use CredentialsTokenProvider instead
type SecretsManagerTokenProvider struct {
	configured bool
	Config     secretsmanagerauth.Config
}

// Ensure SecretsManagerTokenProvider implements TokenProvider interface
var _ TokenProvider = &SecretsManagerTokenProvider{}

func (provider *SecretsManagerTokenProvider) GetRawToken(ctx context.Context) (RawToken, error) {
	if !provider.configured {
		secretsmanagerauth.Configure(provider.Config)
		provider.configured = true
	}

	if err := secretsmanagerauth.SignIn(ctx); err != nil {
		return "", fmt.Errorf("unable to sign-in as service user '%s': %w", provider.Config.SecretKey, err)
	}

	return RawToken(secretsmanagerauth.GetTokens().IdentityToken), nil
}
