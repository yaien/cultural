package models

import (
	"fmt"
	"maps"
	"slices"

	"github.com/yaien/cultural/internal/library/worker"
)

// GetFormat returns the best format for the given variant.
// If the requested width is less than or equal to 0, or if there is only one format available,
// it returns the biggest format. Otherwise, it finds the nearest bigger or equal format based on the requested width.
func (f *File) GetFormat(v int) (format Format, err error) {
	if len(f.Formats) == 0 {
		err = fmt.Errorf("file has no formats")
		return
	}

	variants := slices.Collect(maps.Keys((f.Formats)))
	slices.Sort(variants)

	switch {

	// If the requested width is less than or equal to 0, or if there is only one format available, return the biggest format
	case v <= 0 || len(variants) == 1:
		format = f.Formats[variants[len(variants)-1]]

	// default case: find the near bigger or equal format based on the requested width
	default:
		for index, variant := range variants {

			// Find the first format that is smaller than or equal to the requested width
			if v <= variant {
				format = f.Formats[variant]
				break
			}

			// If we reached the end of the formats and haven't found a suitable one, return the biggest format
			if index == len(variants)-1 {
				format = f.Formats[variant]
				break
			}
		}
	}

	return
}

const GenerateFormatsTaskName = "generate-formats"

// NewGenerateFormatsTask creates a new worker task to generate formats for the file.
func (f *File) NewGenerateFormatsTask() worker.Task {
	return worker.Task{
		Name: GenerateFormatsTaskName,
		Data: map[string]string{"id": f.ID.Hex()},
	}
}
