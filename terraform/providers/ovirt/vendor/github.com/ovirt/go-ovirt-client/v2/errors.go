package ovirtclient

import (
	"errors"
	"fmt"
	"strings"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

// ErrorCode is a code that can be used to identify error types. These errors are identified on a best effort basis
// from the underlying oVirt connection.
type ErrorCode string

// EAccessDenied signals that the provided credentials for the oVirt engine were incorrect.
const EAccessDenied ErrorCode = "access_denied"

// ENotAnOVirtEngine signals that the server did not respond with a proper oVirt response.
const ENotAnOVirtEngine ErrorCode = "not_ovirt_engine"

// ETLSError signals that the provided CA certificate did not match the server that was attempted to connect.
const ETLSError ErrorCode = "tls_error"

// ENotFound signals that the resource requested was not found.
const ENotFound ErrorCode = "not_found"

// EMultipleResults indicates that multiple items were found where only one was expected.
const EMultipleResults ErrorCode = "multiple_results"

// EBug signals an error that should never happen. Please report this.
const EBug ErrorCode = "bug"

// EConnection signals a problem with the connection.
const EConnection ErrorCode = "connection"

// EPermanentHTTPError indicates a HTTP error code that should not be retried.
const EPermanentHTTPError ErrorCode = "permanent_http_error"

// EPending signals that the client library is still waiting for an action to be completed.
const EPending ErrorCode = "pending"

// EUnexpectedDiskStatus indicates that a disk was in a status that was not expected in this state.
const EUnexpectedDiskStatus ErrorCode = "unexpected_disk_status"

// ETimeout signals that the client library has timed out waiting for an action to be completed.
const ETimeout ErrorCode = "timeout"

// EFieldMissing indicates that the oVirt API did not return a specific field. This is most likely a bug, please report
// it.
const EFieldMissing ErrorCode = "field_missing"

// EBadArgument indicates that an input parameter was incorrect.
const EBadArgument ErrorCode = "bad_argument"

// EFileReadFailed indicates that reading a local file failed.
const EFileReadFailed ErrorCode = "file_read_failed"

// EUnexpectedImageTransferPhase indicates that an image transfer was in an unexpected phase.
const EUnexpectedImageTransferPhase ErrorCode = "unexpected_image_transfer_phase"

// EUnidentified is an unidentified oVirt error. When passed to the wrap() function this error code will cause the
// wrap function to look at the wrapped error and either fetch the error code from that error, or identify the error
// from its text.
//
// If you see this error type in a log please report this error so we can add an error code for it.
const EUnidentified ErrorCode = "generic_error"

// EUnsupported signals that an action is not supported. This can indicate a disk format or a combination of parameters.
const EUnsupported ErrorCode = "unsupported"

// EDiskLocked indicates that the disk in question is locked.
const EDiskLocked ErrorCode = "disk_locked"

// EVMLocked indicates that the virtual machine in question is locked.
const EVMLocked ErrorCode = "vm_locked"

// ERelatedOperationInProgress means that the engine is busy working on something else on the same resource.
const ERelatedOperationInProgress ErrorCode = "related_operation_in_progress"

// ELocalIO indicates an input/output error on the client side. For example, a disk could not be read.
const ELocalIO ErrorCode = "local_io_error"

// EConflict indicates an error where you tried to create or update a resource which is already in use in a different,
// conflicting way. For example, you tried to attach a disk that is already attached.
const EConflict ErrorCode = "conflict"

// EHotPlugFailed indicates that a disk could not be hot plugged.
const EHotPlugFailed ErrorCode = "hot_plug_failed"

// EInvalidGrant is an error returned from the oVirt Engine when the SSO token expired. In this case we must reconnect
// and retry the API call.
const EInvalidGrant ErrorCode = "invalid_grant"

// ECannotRunVM indicates an error with the VM configuration which prevents it from being run.
const ECannotRunVM ErrorCode = "cannot_run_vm"

// CanRecover returns true if there is a way to automatically recoverFailure from this error. For the actual recovery an
// appropriate recovery strategy must be passed to the retry function.
func (e ErrorCode) CanRecover() bool {
	switch e {
	case EInvalidGrant:
		return true
	default:
		return false
	}
}

// CanAutoRetry returns false if the given error code is permanent and an automatic retry should not be attempted.
func (e ErrorCode) CanAutoRetry() bool {
	switch e {
	case EBadArgument:
		return false
	case EAccessDenied:
		return false
	case ENotAnOVirtEngine:
		return false
	case ETLSError:
		return false
	case ENotFound:
		return false
	case EMultipleResults:
		return false
	case EBug:
		return false
	case EUnsupported:
		return false
	case EFieldMissing:
		return false
	case EPermanentHTTPError:
		return false
	case EUnexpectedDiskStatus:
		return false
	case ECannotRunVM:
		return false
	default:
		return true
	}
}

// EngineError is an error representation for errors received while interacting with the oVirt engine.
//
// Usage:
//
//   if err != nil {
//     var realErr ovirtclient.EngineError
//     if errors.As(err, &realErr) {
//          // deal with EngineError
//     } else {
//          // deal with other errors
//     }
//   }
type EngineError interface {
	error

	// Message returns the error message without the error code.
	Message() string
	// String returns the string representation for this error.
	String() string
	// HasCode returns true if the current error, or any preceding error has the specified error code.
	HasCode(ErrorCode) bool
	// Code returns an error code for the failure.
	Code() ErrorCode
	// Unwrap returns the underlying error
	Unwrap() error
	// CanRecover indicates that this error can be automatically recovered with the use of the proper recovery strategy
	// passed to the retry function.
	CanRecover() bool
	// CanAutoRetry returns false if an automatic retry should not be attempted.
	CanAutoRetry() bool
}

// HasErrorCode returns true if the specified error has the specified error code.
func HasErrorCode(err error, code ErrorCode) bool {
	var e EngineError
	if errors.As(err, &e) {
		return e.HasCode(code)
	}
	e = realIdentify(err)
	return e.HasCode(code)
}

type engineError struct {
	message string
	code    ErrorCode
	cause   error
}

func (e *engineError) HasCode(code ErrorCode) bool {
	if e.code == code {
		return true
	}
	if cause := e.Unwrap(); cause != nil {
		var causeE EngineError
		if errors.As(cause, &causeE) {
			return causeE.HasCode(code)
		}
	}
	return false
}

func (e *engineError) Message() string {
	return e.message
}

func (e *engineError) String() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *engineError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *engineError) Code() ErrorCode {
	return e.code
}

func (e *engineError) Unwrap() error {
	return e.cause
}

func (e *engineError) CanRecover() bool {
	return e.code.CanRecover()
}

func (e *engineError) CanAutoRetry() bool {
	return e.code.CanAutoRetry()
}

func newFieldNotFound(object string, field string) error {
	return newError(EFieldMissing, "no %s field found on %s object", field, object)
}

func newError(code ErrorCode, format string, args ...interface{}) EngineError {
	return &engineError{
		message: fmt.Sprintf(format, args...),
		code:    code,
	}
}

// wrap wraps an error, adding an error code and message in the process. The wrapped error is added
// to the message automatically in Go style. If the passed error code is EUnidentified or not an EngineError
// this function will attempt to identify the error deeper.
func wrap(err error, code ErrorCode, format string, args ...interface{}) EngineError {
	// gocritic will complain on the following line due to appendAssign, but that's legit here.
	realArgs := append(args, err) // nolint:gocritic
	realMessage := fmt.Sprintf(fmt.Sprintf("%s (%v)", format, "%v"), realArgs...)
	if code == EUnidentified {
		var realErr EngineError
		if errors.As(err, &realErr) {
			code = realErr.Code()
		} else if e := realIdentify(err); e != nil {
			err = e
			code = e.Code()
			realMessage = e.Message()
		}
	}
	return &engineError{
		message: realMessage,
		code:    code,
		cause:   err,
	}
}

func realIdentify(err error) EngineError {
	var authErr *ovirtsdk.AuthError
	var notFoundErr *ovirtsdk.NotFoundError
	switch {
	case strings.Contains(err.Error(), "Cannot run VM without at least one bootable disk."):
		return wrap(
			err,
			ECannotRunVM,
			"cannot run VM due to a missing bootable disk",
		)
	case strings.Contains(err.Error(), "Physical Memory Guaranteed cannot exceed Memory Size"):
		return wrap(
			err,
			EBadArgument,
			"guaranteed memory size must be lower than the memory size",
		)
	case strings.Contains(err.Error(), "stopped after") && strings.Contains(err.Error(), "redirects"):
		return wrap(
			err,
			ENotAnOVirtEngine,
			"the specified oVirt Engine URL has resulted in a redirect, check if your URL is correct",
		)
	case strings.Contains(err.Error(), "parse non-array sso with response"):
		return wrap(
			err,
			ENotAnOVirtEngine,
			"invalid credentials, or the URL does not point to an oVirt Engine, check your settings",
		)
	case strings.Contains(err.Error(), "server gave HTTP response to HTTPS client"):
		return wrap(
			err,
			ENotAnOVirtEngine,
			"the server gave a HTTP response to a HTTPS client, check if your URL is correct",
		)
	case strings.Contains(err.Error(), "invalid_grant: The provided authorization grant for the auth code has expired."):
		return wrap(err, EInvalidGrant, "please reauthenticate for a new access token")
	case strings.Contains(err.Error(), "tls"):
		fallthrough
	case strings.Contains(err.Error(), "x509"):
		return wrap(err, ETLSError, "TLS error, check your CA certificate settings")
	case errors.As(err, &notFoundErr):
		return wrap(err, ENotFound, "the requested resource was not found")
	case strings.Contains(err.Error(), "Disk is locked"):
		return wrap(err, EDiskLocked, "the disk is locked")
	case strings.Contains(err.Error(), "VM is locked"):
		return wrap(err, EVMLocked, "the VM is locked")
	case strings.Contains(err.Error(), "Failed to hot-plug disk"):
		return wrap(err, EHotPlugFailed, "failed to hot-plug disk")
	case strings.Contains(err.Error(), "Related operation is currently in progress."):
		return wrap(err, ERelatedOperationInProgress, "a related operation is in progress")
	case strings.Contains(err.Error(), "Disk configuration") && strings.Contains(err.Error(), " is incompatible with the storage domain type."):
		return wrap(err, EBadArgument, "disk configuration is incompatible with the storage domain type")
	case strings.Contains(err.Error(), "409 Conflict"):
		return wrap(err, EConflict, "conflicting operations")
	case errors.As(err, &authErr):
		fallthrough
	case strings.Contains(err.Error(), "access_denied"):
		return wrap(err, EAccessDenied, "access denied, check your credentials")
	default:
		return nil
	}
}
