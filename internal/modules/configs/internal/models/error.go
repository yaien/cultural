package models

import (
	"errors"
	"net/http"
)

type Error struct {
	Err        error  `json:"-"`
	Code       string `json:"code"`
	HTTPStatus int    `json:"-"`
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func WrapError(err error, code string, httpStatus int) *Error {
	return &Error{
		Err:        err,
		Code:       code,
		HTTPStatus: httpStatus,
	}
}

func DecodeError(err error) *Error {
	return WrapError(err, "decode_error", http.StatusBadRequest)
}

func NotFoundError(err error) *Error {
	return WrapError(err, "not_found", http.StatusNotFound)
}

func IsNotFoundError(err error) bool {
	var e *Error
	if ok := errors.As(err, &e); ok {
		return e.Code == "not_found"
	}
	return false
}
