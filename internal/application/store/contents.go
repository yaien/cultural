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

type Contents struct {
	repository gorm.Interface[Product]
	storage    *storage.Storage
}

func NewContents(db *gorm.DB, storage *storage.Storage) *Contents {
	return &Contents{gorm.G[Product](db), storage}
}

type UploadFileOptions struct {
	PresentationID primitive.UUID
	ProductID      primitive.ID
	OrganizationID primitive.ID
	Name           string
	Size           int64
	ContentType    string
	Data           io.Reader
}

func (c *Contents) Upload(ctx context.Context, req *UploadFileOptions) (*Product, *Presentation, error) {
	product, err := c.repository.
		Where("id = ? AND organization_id = ?", req.ProductID, req.OrganizationID).First(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var presentation *Presentation
	for i := range product.Presentations {
		if product.Presentations[i].ID == req.PresentationID {
			presentation = &product.Presentations[i]
		}
	}

	if presentation == nil {
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.PresentationID)
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
	presentation.Contents = append(presentation.Contents, Content{
		ID:         primitive.NewUUID(),
		FileID:     file.ID,
		FilePreset: file.Preset,
	})

	if _, err = c.repository.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return &product, presentation, nil
}

type GetFileOptions struct {
	ProductID      primitive.ID
	ContentID      primitive.UUID
	PresentationID primitive.UUID
	OrganizationID primitive.ID
}

func (c *Contents) Get(ctx context.Context, req *GetFileOptions) (*Product, *Presentation, *Content, error) {
	product, err := c.repository.
		Where("id = ? AND organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	index := slices.IndexFunc(product.Presentations, func(p Presentation) bool {
		return p.ID == req.PresentationID
	})

	if index == -1 {
		return nil, nil, nil, fmt.Errorf("presentation not found")
	}

	presentation := &product.Presentations[index]

	index = slices.IndexFunc(presentation.Contents, func(c Content) bool {
		return c.ID == req.ContentID
	})

	if index == -1 {
		return nil, nil, nil, fmt.Errorf("content not found")
	}

	content := &presentation.Contents[index]

	return &product, presentation, content, nil
}

type ToggleFilesOptions struct {
	PresentationID primitive.UUID
	ProductID      primitive.ID
	OrganizationID primitive.ID
	ContentIDS     []primitive.UUID
}

func (c *Contents) Toggle(ctx context.Context, req *ToggleFilesOptions) (*Product, *Presentation, error) {
	product, err := c.repository.
		Where("id = ? AND organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	var presentation *Presentation
	if index := slices.IndexFunc(product.Presentations, func(p Presentation) bool { return p.ID == req.PresentationID }); index != -1 {
		presentation = &product.Presentations[index]
	} else {
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.PresentationID)
	}

	if len(presentation.Contents) != len(req.ContentIDS) {
		return nil, nil, coderror.New("presentation_file_count_mismatch", fmt.Errorf("presentation file count mismatch"))
	}

	var contents []Content

	for _, id := range req.ContentIDS {

		if index := slices.IndexFunc(presentation.Contents, func(content Content) bool { return content.ID == id }); index != -1 {
			contents = append(contents, presentation.Contents[index])
			continue
		}

		return nil, nil, coderror.Newf("content_not_found", "content with id %s not found in presentation", id)

	}

	presentation.Contents = contents
	product.UpdatedAt = time.Now()

	if _, err = c.repository.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return &product, presentation, nil
}

type DeleteFileOptions struct {
	PresentationID primitive.UUID
	ContentID      primitive.UUID
	ProductID      primitive.ID
	OrganizationID primitive.ID
}

func (c *Contents) Delete(ctx context.Context, req *DeleteFileOptions) (*Product, *Presentation, error) {
	product, err := c.repository.
		Where("id = ? AND organization_id = ?", req.ProductID, req.OrganizationID).
		First(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("error fetching product: %w", err)
	}

	index := slices.IndexFunc(product.Presentations, func(p Presentation) bool { return p.ID == req.PresentationID })
	if index == -1 {
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.PresentationID)
	}

	presentation := &product.Presentations[index]

	index = slices.IndexFunc(presentation.Contents, func(content Content) bool { return content.ID == req.ContentID })
	if index == -1 {
		return nil, nil, coderror.Newf("content_not_found", "content with id %s not found in presentation", req.ContentID)
	}

	content := presentation.Contents[index]

	err = c.storage.Delete(ctx, req.OrganizationID, content.FileID)
	if err != nil {
		return nil, nil, fmt.Errorf("error deleting file: %w", err)
	}

	presentation.Contents = slices.Delete(presentation.Contents, index, index+1)
	product.UpdatedAt = time.Now()

	if _, err = c.repository.Updates(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product: %w", err)
	}

	return &product, presentation, nil
}
