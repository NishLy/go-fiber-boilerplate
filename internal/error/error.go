package apperror

import (
	"errors"
	"fmt"
)

// Code represents an application-level error code.
type Code int

const (
	NotFound  Code = 404
	Duplicate Code = 409
	Internal  Code = 500
	Invalid   Code = 400
)

// String makes Code implement fmt.Stringer for readable logging.
func (c Code) String() string {
	switch c {
	case NotFound:
		return "NOT_FOUND"
	case Duplicate:
		return "DUPLICATE"
	case Internal:
		return "INTERNAL"
	case Invalid:
		return "INVALID"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(c))
	}
}

// Error is a structured application error carrying a code, human-readable
// message, and an optional underlying cause.
type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	msg := e.Message
	if e.Err != nil {
		msg = e.Err.Error()
	}
	return fmt.Sprintf("[%s] %s", e.Code, msg)
}

func (e *Error) Unwrap() error { return e.Err }

// Is allows errors.Is(err, &Error{Code: NotFound}) matching by code.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// New constructs an Error. Pass nil for err if there is no underlying cause.
func New(code Code, msg string, err error) *Error {
	return &Error{Code: code, Message: msg, Err: err}
}

// Sentinel errors for use with errors.Is.
var (
	ErrNotFound  = &Error{Code: NotFound}
	ErrDuplicate = &Error{Code: Duplicate}
	ErrInternal  = &Error{Code: Internal}
	ErrInvalid   = &Error{Code: Invalid}
)

// Convenience constructors — accept an optional custom message.
// If msg is empty, a default is used.

func NotFoundErr(err error, msg ...string) *Error {
	return New(NotFound, firstOr(msg, "resource not found"), err)
}

func DuplicateErr(err error, msg ...string) *Error {
	return New(Duplicate, firstOr(msg, "duplicate data"), err)
}

func InternalErr(err error, msg ...string) *Error {
	return New(Internal, firstOr(msg, "internal server error"), err)
}

func BadRequestErr(err error, msg ...string) *Error {
	return New(Invalid, firstOr(msg, "bad request"), err)
}

// IsCode reports whether any error in err's chain has the given code.
func IsCode(err error, code Code) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

func firstOr(vals []string, fallback string) string {
	if len(vals) > 0 && vals[0] != "" {
		return vals[0]
	}
	return fallback
}
