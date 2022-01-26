package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrIncorrectCredentials = errors.New("incorrect credentials")
	ErrChallenged           = errors.New("user password needs to be reset")
	ErrTooManyRequests      = errors.New("too many requests to Enlight SSO")
	ErrInactivated          = errors.New("user has been inactivated")
	ErrUnknownTokenType     = errors.New("provided token type not present in response")
)

const (
	DefaultTokenType = "identityToken"
)

type CredentialsTokenProvider struct {
	Username string
	Password string
	Endpoint string

	Client    CredentialsClient
	TokenType string
}

type CredentialsClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Tokens    map[string]RawToken `json:"tokens"`
	Challenge *struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"challenge"`
}

func (provider *CredentialsTokenProvider) GetRawToken(ctx context.Context) (RawToken, error) {
	if provider.TokenType == "" {
		provider.TokenType = DefaultTokenType
	}

	response, err := provider.signIn(ctx, SignInRequest{
		Username: provider.Username,
		Password: provider.Password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign-in: %w", err)
	}

	if response.Challenge != nil {
		return "", ErrChallenged
	}

	token := response.Tokens[provider.TokenType]
	if token == "" {
		return "", fmt.Errorf("%w: %s", ErrUnknownTokenType, provider.TokenType)
	}

	return token, nil
}

func (provider *CredentialsTokenProvider) signIn(ctx context.Context, creds SignInRequest) (*SignInResponse, error) {
	if provider.Client == nil {
		provider.Client = http.DefaultClient
	}

	payload, err := json.Marshal(creds)
	if err != nil {
		return nil, fmt.Errorf("marshalling credentials: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, provider.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("initializing http request: %w", err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	rs, err := provider.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to perform http request: %w", err)
	}

	defer rs.Body.Close()

	var response struct {
		Data  SignInResponse
		Error struct {
			Message string
		}
	}

	if err := json.NewDecoder(rs.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("unmarshalling json response: %w", err)
	}

	switch rs.StatusCode {
	case http.StatusOK:
		return &response.Data, nil
	case http.StatusBadRequest:
		if response.Error.Message == "incorrect username or password" {
			return nil, ErrIncorrectCredentials
		}
	case http.StatusConflict:
		return nil, ErrInactivated
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	}

	return nil, fmt.Errorf("unknown http error %d: %s", rs.StatusCode, response.Error.Message)
}
