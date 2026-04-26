package apperror

type Code string

const (
	NotFound  Code = "NOT_FOUND"
	Duplicate Code = "DUPLICATE"
	Internal  Code = "INTERNAL"
	Invalid   Code = "INVALID"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code Code, msg string, err error) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func NotFoundErr(err error) *Error {
	return New(NotFound, "resource not found", err)
}

func DuplicateErr(err error) *Error {
	return New(Duplicate, "duplicate data", err)
}

func InternalErr(err error) *Error {
	return New(Internal, "internal error", err)
}
