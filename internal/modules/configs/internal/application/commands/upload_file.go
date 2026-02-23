package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/library/worker"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadFileCommand struct {
	storage storage.Storage
	files   models.FileRepository
	queue   *worker.Queue
}

type UploadFileRequest struct {
	Name           string
	Size           int64
	ContentType    string
	Data           io.Reader
	OrganizationID primitive.ObjectID
}

func NewUploadFileCommand(files models.FileRepository, st storage.Storage, q *worker.Queue) *UploadFileCommand {
	return &UploadFileCommand{st, files, q}
}

func (c *UploadFileCommand) UploadFile(ctx context.Context, req *UploadFileRequest) (*models.File, error) {
	_, err := c.files.GetByOrganizationIDAndName(ctx, req.OrganizationID, req.Name)

	var e *models.Error
	switch {
	case err == nil:
		return nil, &models.Error{Code: "name_already_exits", Err: errors.New("file already exists")}
	case errors.As(err, &e) && e.Code == "not_found":
	default:
		return nil, fmt.Errorf("failed to check file existence: %w", err)
	}

	id := primitive.NewObjectID()

	if err = c.storage.Put(id.Hex(), req.Size, req.Data); err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	dir, src, err := c.storage.Mount(id.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to mount file: %w", err)
	}

	defer func() {
		if err := c.storage.Unmount(dir); err != nil {
			slog.Error("Failed to unmount file", "error", err)
		}
	}()

	width, height, variant, err := models.GetFileDimensionByContentType(ctx, src, req.ContentType)
	if err != nil && !errors.Is(err, models.ErrUnsupportedContentType) {
		return nil, fmt.Errorf("failed to get file dimension: %w", err)
	}

	// Extract preset from content type (e.g., "image/jpeg" -> "image")
	preset := strings.Split(req.ContentType, "/")[0]

	// Remove file extension from name (e.g., "photo.jpg" -> "photo")
	name := req.Name
	if idx := strings.LastIndex(req.Name, "."); idx != -1 {
		name = req.Name[:idx-1]
	}

	file := &models.File{
		ID:             id,
		OrganizationID: req.OrganizationID,
		Name:           name,
		Preset:         preset,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Formats: []models.Format{{
			ID:          id,
			Width:       width,
			Height:      height,
			Variant:     variant,
			Size:        req.Size,
			ContentType: req.ContentType,
		}},
	}

	err = c.files.Create(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	if variant > 0 {
		if err = c.queue.Push(ctx, file.NewGenerateFormatsTask()); err != nil {
			return nil, fmt.Errorf("failed to push compress-file job: %w", err)
		}
	}

	return file, nil
}
