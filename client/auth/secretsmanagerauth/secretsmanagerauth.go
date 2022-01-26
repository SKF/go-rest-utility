package secretsmanagerauth

import (
	"context"
	"fmt"

	"github.com/SKF/go-utility/v2/auth"
	"github.com/SKF/go-utility/v2/auth/secretsmanagerauth"
)

type SecretsManagerAuth struct {
	config secretsmanagerauth.Config
}

func New(config secretsmanagerauth.Config) SecretsManagerAuth {
	secretsmanagerauth.Configure(config)
	return SecretsManagerAuth{config: config}
}

func (sma SecretsManagerAuth) GetTokens(ctx context.Context) (auth.Tokens, error) {
	if err := secretsmanagerauth.SignIn(ctx); err != nil {
		return auth.Tokens{}, fmt.Errorf("unable to sign-in using secretKey '%s': %w", sma.config.SecretKey, err)
	}

	return secretsmanagerauth.GetTokens(), nil
}
