package storage

type Error struct {
	code string
	err  error
}

func (e *Error) Code() string {
	return e.code
}
