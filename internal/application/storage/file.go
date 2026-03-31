package storage

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

type File struct {
	ID             primitive.ID `gorm:"primaryKey;autoIncrement"`
	OrganizationID primitive.ID `gorm:"index"`
	Name           string       `gorm:"index"`
	Preset         string
	Formats        []Format `gorm:"type:jsonb;serializer:json"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Format struct {
	ID          string
	Variant     int
	Size        int64
	Width       int
	Height      int
	ContentType string
}

// GetFormat returns the best format for the given variant.
// If the requested width is less than or equal to 0, or if there is only one format available,
// it returns the biggest format. Otherwise, it finds the nearest bigger or equal format based on the requested width.
func (f *File) GetFormat(v int) (format Format, err error) {
	if len(f.Formats) == 0 {
		err = fmt.Errorf("file has no formats")
		return
	}

	// Sort the formats by their variant (width) in ascending order
	slices.SortFunc(f.Formats, func(a, b Format) int { return a.Variant - b.Variant })

	switch {

	// If the requested width is less than or equal to 0, or if there is only one format available, return the biggest format
	case v <= 0 || len(f.Formats) == 1:
		format = f.Formats[len(f.Formats)-1]

	// default case: find the near bigger or equal format based on the requested width
	default:
		for index := range f.Formats {

			// Find the first format that is smaller than or equal to the requested width
			if v <= f.Formats[index].Variant {
				format = f.Formats[index]
				break
			}

			// If we reached the end of the formats and haven't found a suitable one, return the biggest format
			if index == len(f.Formats)-1 {
				format = f.Formats[index]
				break
			}
		}
	}

	return
}
