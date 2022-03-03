package client

import (
	"fmt"
	"io"
	"net/http"
)

var (
	ErrBadRequest                   = newHTTPError(http.StatusBadRequest)
	ErrUnauthorized                 = newHTTPError(http.StatusUnauthorized)
	ErrPaymentRequired              = newHTTPError(http.StatusPaymentRequired)
	ErrForbidden                    = newHTTPError(http.StatusForbidden)
	ErrNotFound                     = newHTTPError(http.StatusNotFound)
	ErrMethodNotAllowed             = newHTTPError(http.StatusMethodNotAllowed)
	ErrNotAcceptable                = newHTTPError(http.StatusNotAcceptable)
	ErrProxyAuthRequired            = newHTTPError(http.StatusProxyAuthRequired)
	ErrRequestTimeout               = newHTTPError(http.StatusRequestTimeout)
	ErrConflict                     = newHTTPError(http.StatusConflict)
	ErrGone                         = newHTTPError(http.StatusGone)
	ErrLengthRequired               = newHTTPError(http.StatusLengthRequired)
	ErrPreconditionFailed           = newHTTPError(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge        = newHTTPError(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong            = newHTTPError(http.StatusRequestURITooLong)
	ErrUnsupportedMediaType         = newHTTPError(http.StatusUnsupportedMediaType)
	ErrRequestedRangeNotSatisfiable = newHTTPError(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed            = newHTTPError(http.StatusExpectationFailed)
	ErrTeapot                       = newHTTPError(http.StatusTeapot)
	ErrMisdirectedRequest           = newHTTPError(http.StatusMisdirectedRequest)
	ErrUnprocessableEntity          = newHTTPError(http.StatusUnprocessableEntity)
	ErrLocked                       = newHTTPError(http.StatusLocked)
	ErrFailedDependency             = newHTTPError(http.StatusFailedDependency)
	ErrTooEarly                     = newHTTPError(http.StatusTooEarly)
	ErrUpgradeRequired              = newHTTPError(http.StatusUpgradeRequired)
	ErrPreconditionRequired         = newHTTPError(http.StatusPreconditionRequired)
	ErrTooManyRequests              = newHTTPError(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge  = newHTTPError(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons   = newHTTPError(http.StatusUnavailableForLegalReasons)

	ErrInternalServerError           = newHTTPError(http.StatusInternalServerError)
	ErrNotImplemented                = newHTTPError(http.StatusNotImplemented)
	ErrBadGateway                    = newHTTPError(http.StatusBadGateway)
	ErrServiceUnavailable            = newHTTPError(http.StatusServiceUnavailable)
	ErrGatewayTimeout                = newHTTPError(http.StatusGatewayTimeout)
	ErrHTTPVersionNotSupported       = newHTTPError(http.StatusHTTPVersionNotSupported)
	ErrVariantAlsoNegotiates         = newHTTPError(http.StatusVariantAlsoNegotiates)
	ErrInsufficientStorage           = newHTTPError(http.StatusInsufficientStorage)
	ErrLoopDetected                  = newHTTPError(http.StatusLoopDetected)
	ErrNotExtended                   = newHTTPError(http.StatusNotExtended)
	ErrNetworkAuthenticationRequired = newHTTPError(http.StatusNetworkAuthenticationRequired)
)

type HTTPError struct {
	StatusCode int
	Status     string

	Instance string
	Body     string
}

func newHTTPError(statusCode int) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
	}
}

func (e HTTPError) withInstance(instance string) HTTPError {
	e.Instance = instance
	return e
}

func (e HTTPError) withBody(reader io.ReadCloser) HTTPError {
	defer reader.Close()

	if body, err := io.ReadAll(reader); err == nil {
		// Assumed to be human readable. Content-type should probably be checked in the future. Or be changed to bytes.
		e.Body = string(body)
	}

	return e
}

func (e HTTPError) Error() string {
	instanceText := ""
	if len(e.Instance) != 0 {
		instanceText = " for : " + e.Instance
	}

	bodyText := e.Body
	if len(bodyText) == 0 {
		bodyText = "[no body]"
	}

	return fmt.Sprintf("got %d%s: %s: %s", e.StatusCode, instanceText, bodyText, e.Status)
}

func (e HTTPError) Is(target error) bool {
	httpErr, ok := target.(HTTPError)
	return ok && httpErr.StatusCode == e.StatusCode
}
