package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetDraftByConfigIDQuery struct {
	drafts models.DraftRepository
}

func NewGetDraftByConfigIDQuery(repo models.DraftRepository) *GetDraftByConfigIDQuery {
	return &GetDraftByConfigIDQuery{drafts: repo}
}

func (q *GetDraftByConfigIDQuery) GetDraftByConfigID(ctx context.Context, configID primitive.ObjectID) (*models.Draft, error) {
	return q.drafts.GetByConfigID(ctx, configID)
}
