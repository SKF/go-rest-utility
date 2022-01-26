package cachedtokenprovider

import (
	"context"
	"fmt"
	"time"

	"github.com/SKF/go-utility/v2/auth/secretsmanagerauth"

	"github.com/SKF/go-rest-utility/client/auth"
)

type Config struct {
	secretsmanagerauth.Config
	TokenTimeToLive time.Duration
}

type Provider struct {
	config Config

	serviceAccountToken string
	lastRefreshTime     time.Time
}

// Ensure Provider implements auth.TokenProvider interface
var _ auth.TokenProvider = Provider{}

func New(config Config) *Provider {
	secretsmanagerauth.Configure(config.Config)

	return &Provider{
		config:          config,
		lastRefreshTime: time.Now(),
	}
}

func (provider Provider) GetRawToken(ctx context.Context) (auth.RawToken, error) {
	if err := provider.refresh(ctx); err != nil {
		return "", fmt.Errorf("failed to refresh identity token: %w", err)
	}

	return auth.RawToken(provider.serviceAccountToken), nil
}

func (provider *Provider) refresh(ctx context.Context) error {
	if provider.tokenIsUninitialized() || provider.tokenIsOutdated() {
		if err := secretsmanagerauth.SignIn(ctx); err != nil {
			return fmt.Errorf("unable to sign-in as service user '%s': %w", provider.config.SecretKey, err)
		}

		tokens := secretsmanagerauth.GetTokens()

		if tokens.IdentityToken == "" {
			return fmt.Errorf("identityToken was empty unexpectedly")
		}

		provider.serviceAccountToken = tokens.IdentityToken
		provider.lastRefreshTime = time.Now()
	}

	return nil
}

func (provider Provider) tokenIsOutdated() bool {
	timeToRefresh := provider.lastRefreshTime.Add(provider.config.TokenTimeToLive)
	return time.Now().After(timeToRefresh)
}

func (provider Provider) tokenIsUninitialized() bool {
	return provider.serviceAccountToken == ""
}
