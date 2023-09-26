package retryafter

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrInvalidFormat = errors.New("invalid header value format")
)

func parseSeconds(retryAfter string) (time.Duration, error) {
	seconds, err := strconv.ParseInt(retryAfter, 10, 64)
	if err != nil {
		return time.Duration(0), err
	}

	if seconds < 0 {
		return time.Duration(0), fmt.Errorf("negative seconds value not allowed")
	}

	return time.Second * time.Duration(seconds), nil
}

func parseHTTPDate(retryAfter string) (time.Duration, error) {
	datetime, err := time.Parse(http.TimeFormat, retryAfter)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Until(datetime), nil
}

func Parse(retryAfter string) (time.Duration, error) {
	if wait, err := parseSeconds(retryAfter); err == nil {
		return wait, nil
	}

	if wait, err := parseHTTPDate(retryAfter); err == nil {
		return wait, nil
	}

	return time.Duration(0), fmt.Errorf("%s: %w", retryAfter, ErrInvalidFormat)
}
