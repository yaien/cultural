package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadFileCommand struct {
	storage storage.Storage
	files   models.FileRepository
}

type UploadFileRequest struct {
	Name           string
	Size           int64
	MimeType       string
	Data           io.Reader
	OrganizationID primitive.ObjectID
}

func NewUploadFileCommand(files models.FileRepository, st storage.Storage) *UploadFileCommand {
	return &UploadFileCommand{st, files}
}

func (c *UploadFileCommand) UploadFile(ctx context.Context, req *UploadFileRequest) (*models.File, error) {
	_, err := c.files.Get(ctx, req.OrganizationID, req.Name)

	var e *models.Error
	switch {
	case err == nil:
		return nil, &models.Error{Code: "name_already_exits", Err: errors.New("file already exists")}
	case errors.As(err, &e) && e.Code == "not_found":
	default:
		return nil, fmt.Errorf("failed to check file existence: %w", err)
	}

	file := &models.File{
		ID:             primitive.NewObjectID(),
		Name:           req.Name,
		Size:           req.Size,
		MimeType:       req.MimeType,
		OrganizationID: req.OrganizationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = c.files.Create(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	err = c.storage.Create(file.ID.Hex(), req.Size, req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	return file, nil
}
