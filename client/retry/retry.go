package retry

import (
	"net/http"
	"time"

	"github.com/SKF/go-rest-utility/internal/retryafter"
)

type Provider interface {
	Should(*http.Request, *http.Response) bool
	Backoff(*http.Response, int) (time.Duration, error)
}

func NewDefaultRetryProvider() Provider {
	return &DefaultRetryProvider{
		backoffProvider: &ExponentialJitterBackoff{
			Base:        1 * time.Millisecond,  //nolint:gomnd
			Cap:         50 * time.Millisecond, //nolint:gomnd
			MaxAttempts: 10,                    //nolint:gomnd
		},
	}
}

type DefaultRetryProvider struct {
	MaxAttempts     int
	backoffProvider BackoffProvider
}

func (d *DefaultRetryProvider) Should(_ *http.Request, resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (d *DefaultRetryProvider) Backoff(resp *http.Response, attempts int) (time.Duration, error) {
	if retryAfter := resp.Header.Get("retry-after"); retryAfter != "" {
		wait, err := retryafter.Parse(retryAfter)
		if err == nil {
			return wait, nil
		}
	}

	return d.backoffProvider.BackoffByAttempt(attempts)
}
