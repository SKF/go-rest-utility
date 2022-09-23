package retry

import (
	"crypto/rand"
	"errors"
	"io"
	"math"
	"math/big"
	"time"
)

type BackoffProvider interface {
	BackoffByAttempt(attempts int) (time.Duration, error)
}

var ErrBackoffExhausted = errors.New("retry attempts exhausted")

const (
	DefaultBackoffBase = time.Nanosecond
	MaxBackoff         = math.MaxInt64
)

// ExponentialJitterBackoff provides an exponentially growing backoff duration
// with jitter to discourage the thundering herd problem.
//
// BackoffByAttempt will return `random_between(0, min(cap, base * 2 ** attempts))`
//
// This is the "Full Jitter" algorithm described at
// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type ExponentialJitterBackoff struct {
	Base        time.Duration // Base backoff duration
	Cap         time.Duration // The maximum backoff time generated
	MaxAttempts int           // If set will throw ErrBackoffExhausted if reached

	JitterSource io.Reader //
}

func (provider *ExponentialJitterBackoff) BackoffByAttempt(attempts int) (time.Duration, error) {
	if provider.JitterSource == nil {
		provider.JitterSource = rand.Reader
	}

	if provider.Base == 0 {
		provider.Base = DefaultBackoffBase
	}

	if provider.MaxAttempts > 0 && attempts > provider.MaxAttempts {
		return 0, ErrBackoffExhausted
	}

	backoff := int64(provider.Base) * int64(1<<attempts)
	if backoff <= 0 {
		backoff = MaxBackoff
	}

	if cap := int64(provider.Cap); cap > 0 && backoff > cap {
		backoff = cap
	}

	jitteredBackoff, err := rand.Int(provider.JitterSource, big.NewInt(backoff))
	if err != nil {
		return 0, err
	}

	return time.Duration(jitteredBackoff.Int64()), nil
}
