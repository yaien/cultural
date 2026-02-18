package queries

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFileQuery struct {
	files   models.FileRepository
	storage storage.Storage
}

func NewGetFileQuery(files models.FileRepository, st storage.Storage) *GetFileQuery {
	return &GetFileQuery{files, st}
}

type GetFileRequest struct {
	OrganizationID primitive.ObjectID
	Name           string
	Quality        int
}

type GetFileResponse struct {
	models.Format
	Name      string
	Data      io.ReadCloser
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *GetFileQuery) GetFile(ctx context.Context, req *GetFileRequest) (*GetFileResponse, error) {
	file, err := q.files.Get(ctx, req.OrganizationID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	format, err := file.GetFormat(req.Quality)
	if err != nil {
		return nil, fmt.Errorf("failed to get file format: %w", err)
	}

	data, err := q.storage.Get(format.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to open file from storage: %w", err)
	}

	res := GetFileResponse{
		Format:    format,
		Name:      file.Name,
		Type:      file.ContentType,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
		Data:      data,
	}

	return &res, nil
}
