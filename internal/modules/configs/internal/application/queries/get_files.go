package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFilesQuery struct {
	files models.FileRepository
}

func NewGetFilesQuery(files models.FileRepository) *GetFilesQuery {
	return &GetFilesQuery{
		files: files,
	}
}

func (q *GetFilesQuery) GetFiles(ctx context.Context, organizationID primitive.ObjectID) ([]*models.File, error) {
	return q.files.List(ctx, organizationID)
}
