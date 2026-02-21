package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	_, err := c.files.Get(ctx, req.OrganizationID, req.Name)

	var e *models.Error
	switch {
	case err == nil:
		return nil, &models.Error{Code: "name_already_exits", Err: errors.New("file already exists")}
	case errors.As(err, &e) && e.Code == "not_found":
	default:
		return nil, fmt.Errorf("failed to check file existence: %w", err)
	}

	id := primitive.NewObjectID()

	err = c.storage.Put(id.Hex(), req.Size, req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	input, err := c.storage.Get(id.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file from storage: %w", err)
	}

	defer input.Close()

	width, height, variant, err := models.GetFileDimension(input, req.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to get file dimensions: %w", err)
	}

	file := &models.File{
		ID:             id,
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Formats: map[int]models.Format{
			variant: {
				ID:          id,
				Width:       width,
				Height:      height,
				Variant:     variant,
				Size:        req.Size,
				ContentType: req.ContentType,
			},
		},
	}

	err = c.files.Create(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	if err = c.queue.Push(ctx, file.NewGenerateFormatsTask()); err != nil {
		return nil, fmt.Errorf("failed to push compress-file job: %w", err)
	}

	return file, nil
}
