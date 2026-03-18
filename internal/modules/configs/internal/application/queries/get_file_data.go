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

type GetFileDataQuery struct {
	files   models.FileRepository
	storage storage.Storage
}

func NewGetFileDataQuery(files models.FileRepository, st storage.Storage) *GetFileDataQuery {
	return &GetFileDataQuery{files, st}
}

type GetFileDataRequest struct {
	OrganizationID primitive.ObjectID
	Name           string
	ID             *primitive.ObjectID
	Variant        int
}

type GetFileDataResponse struct {
	models.Format
	Name      string
	Data      io.ReadCloser
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *GetFileDataQuery) GetFileData(ctx context.Context, req *GetFileDataRequest) (*GetFileDataResponse, error) {
	var file *models.File
	var err error

	if req.ID != nil {
		file, err = q.files.GetByOrganizationIDAndID(ctx, req.OrganizationID, *req.ID)
	} else {
		file, err = q.files.GetByOrganizationIDAndName(ctx, req.OrganizationID, req.Name)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	format, err := file.GetFormat(req.Variant)
	if err != nil {
		return nil, fmt.Errorf("failed to get file format: %w", err)
	}

	data, err := q.storage.Get(format.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to open file from storage: %w", err)
	}

	res := GetFileDataResponse{
		Format:    format,
		Name:      file.Name,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
		Data:      data,
	}

	return &res, nil
}
