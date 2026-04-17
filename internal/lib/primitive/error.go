package primitive

import (
	"errors"

	"github.com/yaien/cultural/internal/lib/coderror"
	"gorm.io/gorm"
)

func Error(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return coderror.New(coderror.NotFound, err)
	default:
		return err
	}
}
