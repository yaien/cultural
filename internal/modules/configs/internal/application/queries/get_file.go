package queries

import (
	"context"
	"fmt"
	"io"

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

func (q *GetFileQuery) GetFile(ctx context.Context, organizationID primitive.ObjectID, name string) (*models.File, io.ReadCloser, error) {
	file, err := q.files.Get(ctx, organizationID, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file: %w", err)
	}

	data, err := q.storage.Get(file.ID.Hex())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file from storage: %w", err)
	}

	return file, data, nil
}
