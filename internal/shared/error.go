package shared

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
