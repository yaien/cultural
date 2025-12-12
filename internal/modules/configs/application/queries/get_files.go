package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFilesQuery struct {
	repository models.FileRepository
}

func NewGetFilesQuery(repository models.FileRepository) *GetFilesQuery {
	return &GetFilesQuery{
		repository: repository,
	}
}

func (q *GetFilesQuery) GetFiles(ctx context.Context, organizationID primitive.ObjectID) ([]*models.File, error) {
	return q.repository.List(ctx, organizationID)
}
