package weberr

import (
	goerrors "errors"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// NoType error - Default placeholder for un-typed errors.
	// Methods on NoType behave differently than other types - they
	// will preserve an existing type when wrapping a typed error.
	NoType ErrorType = iota

	// 4xx Client errors
	// -----------------

	// BadRequest error - Code 400
	BadRequest ErrorType = http.StatusBadRequest
	// Unauthorized error - Code 401
	Unauthorized ErrorType = http.StatusUnauthorized
	// PaymentRequired error - Code 402
	PaymentRequired ErrorType = http.StatusPaymentRequired
	// Forbidden error - Code 403
	Forbidden ErrorType = http.StatusForbidden
	// NotFound error - Code 404
	NotFound ErrorType = http.StatusNotFound
	// MethodNotAllowed error - Code 405
	MethodNotAllowed ErrorType = http.StatusMethodNotAllowed
	// NotAcceptable error - Code 406
	NotAcceptable ErrorType = http.StatusNotAcceptable
	// ProxyAuthRequired error - Code 407
	ProxyAuthRequired ErrorType = http.StatusProxyAuthRequired
	// RequestTimeout error - Code 408
	RequestTimeout ErrorType = http.StatusRequestTimeout
	// Conflict error - Code 409
	Conflict ErrorType = http.StatusConflict
	// Gone error - Code 410
	Gone ErrorType = http.StatusGone
	// LengthRequired error - Code 411
	LengthRequired ErrorType = http.StatusLengthRequired
	// PreconditionFailed error - Code 412
	PreconditionFailed ErrorType = http.StatusPreconditionFailed
	// RequestEntityTooLarge error - Code 413
	RequestEntityTooLarge ErrorType = http.StatusRequestEntityTooLarge
	// RequestURITooLong error - Code 414
	RequestURITooLong ErrorType = http.StatusRequestURITooLong
	// UnsupportedMediaType error - Code 415
	UnsupportedMediaType ErrorType = http.StatusUnsupportedMediaType
	// RequestedRangeNotSatisfiable error - Code 416
	RequestedRangeNotSatisfiable ErrorType = http.StatusRequestedRangeNotSatisfiable
	// ExpectationFailed error - Code 417
	ExpectationFailed ErrorType = http.StatusExpectationFailed
	// Teapot error - Code 418
	Teapot ErrorType = http.StatusTeapot
	// UnprocessableEntity error - Code 422
	UnprocessableEntity ErrorType = http.StatusUnprocessableEntity
	// Locked error - Code 423
	Locked ErrorType = http.StatusLocked
	// FailedDependency error - Code 424
	FailedDependency ErrorType = http.StatusFailedDependency
	// UpgradeRequired error - Code 426
	UpgradeRequired ErrorType = http.StatusUpgradeRequired
	// PreconditionRequired error - Code 428
	PreconditionRequired ErrorType = http.StatusPreconditionRequired
	// TooManyRequests error - Code 429
	TooManyRequests ErrorType = http.StatusTooManyRequests
	// RequestHeaderFieldsTooLarge error - Code 431
	RequestHeaderFieldsTooLarge ErrorType = http.StatusRequestHeaderFieldsTooLarge
	// UnavailableForLegalReasons error - Code 451
	UnavailableForLegalReasons ErrorType = http.StatusUnavailableForLegalReasons

	// 5xx Server errors
	// -----------------

	// InternalServerError error - Code 500
	InternalServerError ErrorType = http.StatusInternalServerError
	// NotImplemented error - Code 501
	NotImplemented ErrorType = http.StatusNotImplemented
	// BadGateway error - Code 502
	BadGateway ErrorType = http.StatusBadGateway
	// ServiceUnavailable error - Code 503
	ServiceUnavailable ErrorType = http.StatusServiceUnavailable
	// GatewayTimeout error - Code 504
	GatewayTimeout ErrorType = http.StatusGatewayTimeout
	// HTTPVersionNotSupported error - Code 505
	HTTPVersionNotSupported ErrorType = http.StatusHTTPVersionNotSupported
	// VariantAlsoNegotiates error - Code 506
	VariantAlsoNegotiates ErrorType = http.StatusVariantAlsoNegotiates
	// InsufficientStorage error - Code 507
	InsufficientStorage ErrorType = http.StatusInsufficientStorage
	// LoopDetected error - Code 508
	LoopDetected ErrorType = http.StatusLoopDetected
	// NotExtended error - Code 510
	NotExtended ErrorType = http.StatusNotExtended
	// NetworkAuthenticationRequired error - Code 511
	NetworkAuthenticationRequired ErrorType = http.StatusNetworkAuthenticationRequired
)

// customError wraps an error with type, separate user message, and details.
type customError struct {
	error
	errorType   ErrorType
	userMessage string
	details     []interface{}
}

// causer interface allows unwrapping an error.
// causer is also used in github.com/pkg/errors
type causer interface {
	Cause() error
}

// Cause unwraps error
func (c *customError) Cause() error { return c.error }

// typed interface identifies error with a type
type typed interface {
	Type() ErrorType
}

// Type returns the error type
func (c *customError) Type() ErrorType { return c.errorType }

// GetType returns the error type for all errors.
// If error is not `typed` - it returns NoType.
func GetType(err error) ErrorType {
	if typeErr, ok := err.(typed); ok {
		return typeErr.Type()
	}

	return NoType
}

// userMessager identifies an error with a user message
type userMessager interface {
	UserMessage() string
}

// UserMessage returns the user message
func (c *customError) UserMessage() string { return c.userMessage }

// GetUserMessage returns user readable error message for all errors.
// If error is not `userMessager` returns empty string.
func GetUserMessage(err error) string {
	if msgErr, ok := err.(userMessager); ok {
		return msgErr.UserMessage()
	}

	return ""
}

// errorDetailer identifies an error with details
type errorDetailer interface {
	Details() []interface{}
}

// Details returns the error details
func (c *customError) Details() []interface{} { return c.details }

// GetDetails returns a slice of arbitrary details for all errors.
// If error is not `errorDetailer` returns nil.
func GetDetails(err error) []interface{} {
	if detailedError, ok := err.(errorDetailer); ok {
		return detailedError.Details()
	}

	return nil
}

// Errorf creates a new error of this type with formatted string.
func (errorType ErrorType) Errorf(msg string, args ...interface{}) error {
	return &customError{
		error:     errors.WithStack(errors.Errorf(msg, args...)),
		errorType: errorType,
	}
}

// Wrapf creates a wrapping error of this type, with formatted string.
// The relation to the wrapped err is implicit, do not add a %s for it (like you would with fmt.Errorf).
// If wrapped err is nil, still returns a new error.
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	if err == nil {
		return errorType.Errorf(msg, args...)
	}

	c := new(customError)
	c.error = errors.Wrapf(err, msg, args...)
	c.userMessage = GetUserMessage(err)
	c.details = GetDetails(err)

	if errorType != NoType {
		c.errorType = errorType
	} else {
		c.errorType = GetType(err)
	}

	return c
}

// UserWrapf adds a formatted user readable message to an error.
// If wrapped err already had a user message, combines them with a colon, you should not add a %s for it.
// Also sets error type (or preserves existing type if called on NoType).
// If wrapped err is nil, still returns a new error.
func (errorType ErrorType) UserWrapf(err error, msg string, args ...interface{}) error {
	if err == nil {
		return errorType.UserErrorf(msg, args...)
	}

	userMsg := fmt.Sprintf(msg, args...)

	c := new(customError)
	c.error = errors.WithStack(err)
	c.details = GetDetails(err)

	origMsg := GetUserMessage(err)
	if origMsg != "" {
		userMsg = fmt.Sprintf("%s: %s", userMsg, origMsg)
	}
	c.userMessage = userMsg

	if errorType != NoType {
		c.errorType = errorType
	} else {
		c.errorType = GetType(err)
	}

	return c
}

// UserErrorf creates a new error with a user readable message.
func (errorType ErrorType) UserErrorf(msg string, args ...interface{}) error {
	message := fmt.Sprintf(msg, args...)
	return &customError{
		error:       errors.WithStack(errors.New(message)),
		errorType:   errorType,
		userMessage: message,
	}
}

// AddDetails adds a details element to an error.
// Also sets error type (or preserves existing type if called on NoType).
func (errorType ErrorType) AddDetails(err error, details interface{}) error {
	if details == nil {
		return err
	}

	if err == nil {
		return errorType.details(details)
	}

	c := new(customError)
	c.error = errors.WithStack(err)
	c.userMessage = GetUserMessage(err)

	c.details = append(GetDetails(err), details)

	if errorType != NoType {
		c.errorType = errorType
	} else {
		c.errorType = GetType(err)
	}

	return c
}

// details creates a new error with arbitrary details
func (errorType ErrorType) details(details interface{}) error {
	return &customError{
		error:     errors.WithStack(errors.New("")),
		errorType: errorType,
		details:   []interface{}{details},
	}
}

// Set the type of the error
func (errorType ErrorType) Set(err error) error {
	if err == nil {
		return nil
	}

	return &customError{
		error:       errors.WithStack(err),
		errorType:   errorType,
		userMessage: GetUserMessage(err),
		details:     GetDetails(err),
	}
}

func (errorType ErrorType) SetUserMessage(err error, msg string) error {
	if err == nil {
		return nil
	}

	var newType ErrorType
	if errorType != NoType {
		newType = errorType
	} else {
		newType = GetType(err)
	}

	return &customError{
		error:       errors.WithStack(err),
		errorType:   newType,
		userMessage: msg,
		details:     GetDetails(err),
	}
}

// Errorf returns a new NoType error with formatted string.
func Errorf(msg string, args ...interface{}) error {
	return NoType.Errorf(msg, args...)
}

// Wrapf creates a wrapping error, with unmodified type and formatted string.
// The relation to the wrapped err is implicit, do not add a %s for it (like you would with fmt.Errorf).
// If wrapped err is nil, still returns a new error.
func Wrapf(err error, msg string, args ...interface{}) error {
	return NoType.Wrapf(err, msg, args...)
}

// UserErrorf returns an error with formatted user message.
func UserErrorf(msg string, args ...interface{}) error {
	return NoType.UserErrorf(msg, args...)
}

// UserWrapf adds a user readable message to an error.
func UserWrapf(err error, msg string, args ...interface{}) error {
	return NoType.UserWrapf(err, msg, args...)
}

// AddDetails adds arbitrary details to an error.
func AddDetails(err error, details interface{}) error {
	return NoType.AddDetails(err, details)
}

// SetUserMessage sets the user message to an error
func SetUserMessage(err error, msg string) error {
	return NoType.SetUserMessage(err, msg)
}

// stackTracer interface is internally defined in github.com/pkg/errors
// and identifies an error with a stack trace
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// baseStackTracer is a helper function to allow reaching
// the initial wrapper that has a stack trace
func baseStackTracer(err error) error {

	if cause, ok := err.(causer); ok {
		candidate := baseStackTracer(cause.Cause())
		if candidate != nil {
			return candidate
		}

		if _, ok := err.(stackTracer); ok {
			return err
		}
	}
	return nil
}

// GetStackTrace returns the stack trace starting from the first error
// that has been wrapped / created
func GetStackTrace(err error) string {
	if err == nil {
		return ""
	}

	err = baseStackTracer(err)
	x, ok := err.(stackTracer)
	if !ok {
		// The error doesn't have a stack trace attached to it
		return fmt.Sprintf("%+v", err)
	}

	st := x.StackTrace()
	return fmt.Sprintf("%+v\n", st[1:])
}

// As finds the first error in err's chain that matches target,
// and if so, sets target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	return goerrors.As(err, target)
}

// Is finds the first error in err's chain that matches target,
// and if so, sets target to that error value and returns true. Otherwise, it returns false.
func Is(err error, target error) bool {
	return goerrors.Is(err, target)
}
