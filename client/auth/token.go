package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type RawToken string

// ParseExpires attempts to extract the `exp` claim from the JWT. This function
// does not validate that the token is signed and should only be used if the
// correctness of this value is not used for security.
func (token RawToken) ParseExpires() (time.Time, error) {
	parts := strings.Split(string(token), ".")
	if len(parts) < 3 { //nolint:gomnd // A JWT should contain 3 parts divided by .
		return time.Time{}, fmt.Errorf("%w: missing parts, found %d should be 3", ErrInvalidToken, len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: not base64 decodeable: %s", ErrInvalidToken, err)
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}

	if err = json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, err
	}

	return time.Unix(claims.Exp, 0), nil
}

func (token RawToken) GetRawToken(ctx context.Context) (RawToken, error) {
	return token, nil
}

func (token RawToken) String() string {
	return string(token)
}
