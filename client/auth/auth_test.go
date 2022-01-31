package auth_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SKF/go-utility/v2/uuid"
	"github.com/stretchr/testify/require"

	"github.com/SKF/go-rest-utility/client/auth"
)

type TestAccessToken struct {
	Email     string
	Lifetime  time.Duration
	IssueTime time.Time
}

func (token TestAccessToken) Build(t *testing.T) auth.RawToken {
	t.Helper()

	if token.IssueTime.IsZero() {
		token.IssueTime = time.Now()
	}

	header, err := encodeToBase64(map[string]string{
		"alg": "none",
	})
	require.NoError(t, err)

	payload, err := encodeToBase64(map[string]interface{}{
		"sub":       uuid.New(),
		"event_id":  uuid.New(),
		"token_use": "access",
		"scope":     "aws.cognito.signin.user.admin",
		"auth_time": token.IssueTime.Unix(),
		"iss":       "https://example.com",
		"exp":       token.IssueTime.Add(token.Lifetime).Unix(),
		"iat":       token.IssueTime.Unix(),
		"jti":       uuid.New(),
		"client_id": "",
		"username":  token.Email,
	})
	require.NoError(t, err)

	return auth.RawToken(header + "." + payload + ".")
}

func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.RawURLEncoding, &buf)

	if err := json.NewEncoder(encoder).Encode(v); err != nil {
		return "", err
	}

	encoder.Close()

	return buf.String(), nil
}

type SSO struct {
	*httptest.Server

	users map[string]*TestUser
}

type TestUser struct {
	Username string
	Password string

	Inactive    bool
	Challenged  bool
	Ratelimited bool

	validTokens map[auth.RawToken]bool
}

func NewSSO() *SSO {
	sso := &SSO{
		users: make(map[string]*TestUser),
	}

	sso.Server = httptest.NewServer(sso)

	return sso
}

func (api *SSO) WithUser(user TestUser) *SSO {
	user.validTokens = make(map[auth.RawToken]bool)

	api.users[user.Username] = &user

	return api
}

func (api *SSO) RequireValidToken(t *testing.T, user string, token auth.RawToken) {
	t.Helper()

	require.Contains(t, api.users, user)
	require.Contains(t, api.users[user].validTokens, token)
}

func (api *SSO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request auth.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":{"message":"unable to decode request"}}`)

		return
	}

	user, found := api.users[request.Username]
	if !found || user.Password != request.Password {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"message":"incorrect username or password"}}`)

		return
	}

	if user.Ratelimited {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"error":{"message":"Too many requests"}}`)

		return
	}

	if user.Inactive {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{"error":{"message":"user is in status inactive"}}`)

		return
	}

	if user.Challenged {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"challenge":{"id":"very-long-string","type":"NEW_PASSWORD_REQUIRED"}}}`)

		return
	}

	identityToken := TestAccessToken{
		Email:    user.Username,
		Lifetime: time.Hour,
	}.Build(new(testing.T))

	user.validTokens[identityToken] = true

	response := struct {
		Data auth.SignInResponse
	}{
		Data: auth.SignInResponse{
			Tokens: map[string]auth.RawToken{
				"identityToken": identityToken,
			},
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error":{"message":"unable to encode response"}}`)

		return
	}
}
