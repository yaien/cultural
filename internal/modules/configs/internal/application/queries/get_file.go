package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFileQuery struct {
	files models.FileRepository
}

func NewGetFileQuery(files models.FileRepository) *GetFileQuery {
	return &GetFileQuery{files}
}

func (q *GetFileQuery) GetFile(ctx context.Context, organizationID primitive.ObjectID, name string) (*models.File, error) {
	return q.files.GetByOrganizationIDAndName(ctx, organizationID, name)
}
