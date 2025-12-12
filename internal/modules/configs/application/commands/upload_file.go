package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadFileCommand struct {
	storage    storage.Storage
	repository models.FileRepository
}

type UploadFileCommandRequest struct {
	Name           string
	Size           int64
	Data           io.Reader
	OrganizationID primitive.ObjectID
}

func NewUploadFileCommand(repo models.FileRepository, st storage.Storage) *UploadFileCommand {
	return &UploadFileCommand{st, repo}
}

func (c *UploadFileCommand) UploadFile(ctx context.Context, req *UploadFileCommandRequest) (*models.File, error) {
	_, err := c.repository.Get(ctx, req.OrganizationID, req.Name)

	var e *models.Error
	switch {
	case err == nil:
		return nil, &models.Error{Code: "name_already_exits", Err: errors.New("file already exists")}
	case errors.As(err, &e) && e.Code == "file_not_found":
	default:
		return nil, fmt.Errorf("failed to check file existence: %w", err)
	}

	data, err := io.ReadAll(req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	file := &models.File{
		ID:             primitive.NewObjectID(),
		Name:           req.Name,
		Size:           req.Size,
		MimeType:       http.DetectContentType(data),
		OrganizationID: req.OrganizationID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = c.repository.Create(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	err = c.storage.Create(file.ID.Hex(), req.Size, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	return file, nil
}
