package coderror

import (
	"errors"
	"fmt"
	"net/http"
)

const NotFound = "not_found"
const DecodeFailed = "decode_failed"

type Error struct {
	Code string
	err  error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func New(code string, err error) *Error {
	return &Error{
		Code: code,
		err:  err,
	}
}

func Newf(code string, format string, args ...any) *Error {
	return &Error{
		Code: code,
		err:  fmt.Errorf(format, args...),
	}
}

func Is(err error, code string) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

func (e *Error) HTTPStatus() int {
	switch e.Code {
	case NotFound:
		return http.StatusNotFound
	case DecodeFailed:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
