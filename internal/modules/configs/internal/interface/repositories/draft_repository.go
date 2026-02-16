package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.DraftRepository = (*DraftRepository)(nil)

type DraftRepository struct {
	db *mongo.Database
}

func NewDraftRepository(db *mongo.Database) *DraftRepository {
	return &DraftRepository{db: db}
}

func (r *DraftRepository) Update(ctx context.Context, draft *models.Draft) error {
	_, err := r.db.Collection("drafts").UpdateOne(ctx, map[string]any{"_id": draft.ID}, map[string]any{"$set": draft})
	return err
}

func (r *DraftRepository) GetByConfigID(ctx context.Context, configID primitive.ObjectID) (*models.Draft, error) {
	var draft models.Draft
	err := r.db.Collection("drafts").FindOne(ctx, map[string]any{"configId": configID}).Decode(&draft)
	return &draft, translate(err)
}
