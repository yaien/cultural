package store

import (
	"context"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"

	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/coderror"
)

var MaxFilesPerPresentation = 5

type Files struct {
	repository gorm.Interface[Product]
	storage    *storage.Storage
}

func NewFiles(db *gorm.DB, storage *storage.Storage) *Files {
	return &Files{gorm.G[Product](db), storage}
}

type UploadFileOptions struct {
	PresentationID primitive.ID
	ProductID      primitive.ID
	OrganizationID primitive.ID
	Name           string
	Size           int64
	ContentType    string
	Data           io.Reader
}

func (c *Files) Upload(ctx context.Context, req *UploadFileOptions) (*Product, *Presentation, error) {
	product, err := c.repository.
		Preload("Presentations.Contents.File", nil).
		Where("id = ? AND organization_id = ?", req.ProductID, req.OrganizationID).First(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var presentation *Presentation
	for _, p := range product.Presentations {
		if p.ID == req.PresentationID {
			presentation = p
		}
	}

	if presentation == nil {
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %d not found", req.PresentationID)
	}

	if len(presentation.Contents) >= MaxFilesPerPresentation {
		return nil, nil, coderror.New("presentation_file_limit_exceeded", fmt.Errorf("presentation file limit exceeded"))
	}

	file, err := c.storage.Upload(ctx, &storage.UploadOptions{
		Name:           req.Name,
		Size:           req.Size,
		ContentType:    req.ContentType,
		Data:           req.Data,
		OrganizationID: req.OrganizationID,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error uploading file: %w", err)
	}

	product.UpdatedAt = time.Now()
	presentation.Contents = append(presentation.Contents, &Content{
		PresentationID: presentation.ID,
		FileID:         file.ID,
		Order:          len(presentation.Contents),
	})

	if _, err = c.repository.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return &product, presentation, nil
}

type ToggleFilesOptions struct {
	PresentationID primitive.ID
	ProductID      primitive.ID
	OrganizationID primitive.ID
	ContentIDS     []primitive.ID
}

func (c *Files) Toggle(ctx context.Context, req *ToggleFilesOptions) (*Product, *Presentation, error) {
	product, err := c.repository.
		Preload("Presentations.Contents.File", nil).
		Where("id = ? AND organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var presentation *Presentation
	for _, p := range product.Presentations {
		if p.ID == req.PresentationID {
			presentation = p
			break
		}
	}

	if presentation == nil {
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %d not found", req.PresentationID)
	}

	if len(presentation.Contents) != len(req.ContentIDS) {
		return nil, nil, coderror.New("presentation_file_count_mismatch", fmt.Errorf("presentation file count mismatch"))
	}

	for i, id := range req.ContentIDS {

		index := slices.IndexFunc(presentation.Contents, func(content *Content) bool { return content.ID.Equal(id) })
		if index == -1 {
			return nil, nil, coderror.Newf("file_not_found", "file with id %d not found in presentation", id)
		}

		presentation.Contents[index].Order = i
	}

	product.UpdatedAt = time.Now()

	if _, err = c.repository.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return &product, presentation, nil
}
