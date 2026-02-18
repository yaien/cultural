package models

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/yaien/cultural/internal/library/worker"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	OrganizationID primitive.ObjectID `bson:"organizationId" json:"organizationId"`
	Name           string             `bson:"name" json:"name"`
	ContentType    string             `bson:"contentType" json:"contentType"`
	Formats        map[int]Format     `bson:"formats" json:"formats"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Format struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Quality int                `bson:"quality" json:"quality"`
	Size    int64              `bson:"size" json:"size"`
	Width   int                `bson:"width" json:"width"`
	Height  int                `bson:"height" json:"height"`
}

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	Get(ctx context.Context, organizationID primitive.ObjectID, name string) (*File, error)
	Rename(ctx context.Context, organizationID primitive.ObjectID, oldName, newName string) error
	Delete(ctx context.Context, organizationID primitive.ObjectID, name string) error
	List(ctx context.Context, organizationID primitive.ObjectID) ([]*File, error)
}

type FileURLFunc func(filename string) string

func NewExternalFileURLFunc(serverURL string, organizationID primitive.ObjectID) FileURLFunc {
	return func(filename string) string {
		return fmt.Sprintf("%s/assets/external/%s/%s", serverURL, organizationID.Hex(), filename)
	}
}

// GetFormat returns the best format for the given quality (width).
// If the requested width is less than or equal to 0, or if there is only one format available,
// it returns the biggest format. Otherwise, it finds the nearest bigger or equal format based on the requested width.
func (f *File) GetFormat(q int) (format Format, err error) {
	if len(f.Formats) == 0 {
		err = fmt.Errorf("file has no formats")
		return
	}

	qualities := slices.Collect(maps.Keys((f.Formats)))
	slices.Sort(qualities)

	switch {

	// If the requested width is less than or equal to 0, or if there is only one format available, return the biggest format
	case q <= 0 || len(qualities) == 1:
		format = f.Formats[qualities[len(qualities)-1]]

	// default case: find the near bigger or equal format based on the requested width
	default:
		for index, quality := range qualities {

			// Find the first format that is smaller than or equal to the requested width
			if q <= quality {
				format = f.Formats[quality]
				break
			}

			// If we reached the end of the formats and haven't found a suitable one, return the biggest format
			if index == len(qualities)-1 {
				format = f.Formats[quality]
				break
			}
		}
	}

	return
}

// NewGenerateFormatsTask creates a new worker task to generate formats for the file.
func (f *File) NewGenerateFormatsTask() worker.Task {
	return worker.Task{
		Name: "generate-formats",
		Data: map[string]string{"id": f.ID.Hex()},
	}
}
