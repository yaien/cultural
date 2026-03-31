package store

import (
	"context"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MaxFilesPerPresentation = 5

type Files struct {
	repository Repository
	storage    *storage.Storage
}

func NewFiles(repository Repository, storage *storage.Storage) *Files {
	return &Files{repository: repository, storage: storage}
}

type UploadFileOptions struct {
	PresentationID primitive.ObjectID
	ProductID      primitive.ObjectID
	OrganizationID primitive.ObjectID
	Name           string
	Size           int64
	ContentType    string
	Data           io.Reader
}

func (c *Files) Upload(ctx context.Context, req *UploadFileOptions) (*Product, *Presentation, error) {
	product, err := c.repository.GetByIDAndOrganizationID(ctx, req.ProductID, req.OrganizationID)
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
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.PresentationID.Hex())
	}

	if len(presentation.Files) >= MaxFilesPerPresentation {
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
	presentation.Files = append(presentation.Files, &File{
		ID:     file.ID,
		Preset: file.Preset,
	})

	if err = c.repository.Update(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return product, presentation, nil
}

type ToggleFilesOptions struct {
	PresentationID primitive.ObjectID
	ProductID      primitive.ObjectID
	OrganizationID primitive.ObjectID
	FileIDS        []primitive.ObjectID
}

func (c *Files) Toggle(ctx context.Context, req *ToggleFilesOptions) (*Product, *Presentation, error) {
	product, err := c.repository.GetByIDAndOrganizationID(ctx, req.ProductID, req.OrganizationID)
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
		return nil, nil, coderror.Newf("presentation_not_found", "presentation with id %s not found", req.PresentationID.Hex())
	}

	if len(presentation.Files) != len(req.FileIDS) {
		return nil, nil, coderror.New("presentation_file_count_mismatch", fmt.Errorf("presentation file count mismatch"))
	}

	var files []*File

	for _, id := range req.FileIDS {
		index := slices.IndexFunc(presentation.Files, func(file *File) bool { return file.ID == id })
		if index == -1 {
			return nil, nil, coderror.Newf("file_not_found", "file with id %s not found in presentation", id.Hex())
		}

		files = append(files, presentation.Files[index])

	}

	presentation.Files = files
	product.UpdatedAt = time.Now()

	if err = c.repository.Update(ctx, product); err != nil {
		return nil, nil, fmt.Errorf("error updating product presentation files: %w", err)
	}

	return product, presentation, nil
}
