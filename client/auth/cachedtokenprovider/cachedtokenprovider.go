package cachedtokenprovider

import (
	"context"
	"fmt"
	"time"

	auth_model "github.com/SKF/go-utility/v2/auth"
	"github.com/SKF/go-utility/v2/auth/secretsmanagerauth"

	"github.com/SKF/go-rest-utility/client/auth"
	auth_client "github.com/SKF/go-rest-utility/client/auth/secretsmanagerauth"
)

type SecretsManagerAuth interface {
	GetTokens(ctx context.Context) (auth_model.Tokens, error)
}

type Config struct {
	secretsmanagerauth.Config
	TokenTimeToLive time.Duration
}

type Provider struct {
	secretsManagerAuth SecretsManagerAuth
	config             Config

	token           string
	lastRefreshTime time.Time
}

// Ensure Provider implements auth.TokenProvider interface
var _ auth.TokenProvider = &Provider{}

// New initializes a Provider and configures the default secrets manager auth implementation.
func New(config Config) *Provider {
	return &Provider{
		secretsManagerAuth: auth_client.New(config.Config),
		config:             config,

		lastRefreshTime: time.Now(),
	}
}

func NewWithCustomAuth(config Config, secretsManagerAuth SecretsManagerAuth) *Provider {
	return &Provider{
		secretsManagerAuth: secretsManagerAuth,
		config:             config,

		lastRefreshTime: time.Now(),
	}
}

func (provider *Provider) GetRawToken(ctx context.Context) (auth.RawToken, error) {
	if err := provider.refresh(ctx); err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return auth.RawToken(provider.token), nil
}

func (provider *Provider) refresh(ctx context.Context) error {
	if provider.tokenIsUninitialized() || provider.tokenIsOutdated() {
		tokens, err := provider.secretsManagerAuth.GetTokens(ctx)
		if err != nil {
			return fmt.Errorf("failed to get tokens from secrets manager auth: %w", err)
		}

		if tokens.IdentityToken == "" {
			return fmt.Errorf("identityToken was empty unexpectedly")
		}

		provider.token = tokens.IdentityToken
		provider.lastRefreshTime = time.Now()
	}

	return nil
}

func (provider Provider) tokenIsOutdated() bool {
	timeToRefresh := provider.lastRefreshTime.Add(provider.config.TokenTimeToLive)
	return time.Now().After(timeToRefresh)
}

func (provider Provider) tokenIsUninitialized() bool {
	return provider.token == ""
}
