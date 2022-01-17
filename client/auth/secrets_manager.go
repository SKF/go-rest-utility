package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sm_v2 "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	sm_v1 "github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretCredentialsTokenProvider struct {
	SecretID      string
	SecretsClient SecretsClient

	Client    CredentialsClient
	TokenType string

	sourceProvider *CredentialsTokenProvider
}

type SecretCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"signInUrl"`
}

type SecretsClient interface {
	GetSecretByID(context.Context, string) ([]byte, error)
}

type SecretsManagerV1Client struct {
	sm_v1.SecretsManager
}

type SecretsManagerV2Client struct {
	*sm_v2.Client
}

func (c SecretsManagerV1Client) GetSecretByID(ctx context.Context, secretID string) ([]byte, error) {
	output, err := c.GetSecretValueWithContext(ctx, &sm_v1.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return nil, err
	}

	return output.SecretBinary, nil
}

func (c SecretsManagerV2Client) GetSecretByID(ctx context.Context, secretID string) ([]byte, error) {
	output, err := c.GetSecretValue(ctx, &sm_v2.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return nil, err
	}

	return output.SecretBinary, nil
}

func (provider *SecretCredentialsTokenProvider) GetRawToken(ctx context.Context) (_ RawToken, err error) {
	generatedCredentials := false

	if provider.sourceProvider == nil {
		provider.sourceProvider, err = provider.generateSourceProvider(ctx)
		if err != nil {
			return "", fmt.Errorf("retreiving secret credentials: %w", err)
		}

		generatedCredentials = true
	}

	token, err := provider.sourceProvider.GetRawToken(ctx)

	// If the sourceProvider was cached and they seem out of date, attempt to regenerate them and try again.
	if errors.Is(err, ErrIncorrectCredentials) && !generatedCredentials {
		provider.sourceProvider = nil
		return provider.GetRawToken(ctx)
	}

	return token, err
}

func (provider *SecretCredentialsTokenProvider) generateSourceProvider(ctx context.Context) (*CredentialsTokenProvider, error) {
	output, err := provider.SecretsClient.GetSecretByID(ctx, provider.SecretID)
	if err != nil {
		return nil, err
	}

	var secret SecretCredentials
	if err := json.Unmarshal(output, &secret); err != nil {
		return nil, fmt.Errorf("credentials not decodable: %w", err)
	}

	return &CredentialsTokenProvider{
		Username:  secret.Username,
		Password:  secret.Password,
		Endpoint:  secret.Endpoint,
		Client:    provider.Client,
		TokenType: provider.TokenType,
	}, nil
}
