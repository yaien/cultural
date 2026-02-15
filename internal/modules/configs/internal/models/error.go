package models

import "errors"

type Error struct {
	Err  error  `json:"-"`
	Code string `json:"code"`
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func WrapError(err error, code string) *Error {
	return &Error{
		Err:  err,
		Code: code,
	}
}

func NotFoundError(err error) *Error {
	return WrapError(err, "not_found")
}

func IsNotFoundError(err error) bool {
	var e *Error
	if ok := errors.As(err, &e); ok {
		return e.Code == "not_found"
	}
	return false
}
