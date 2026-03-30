package label

import (
	"context"

	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ DraftRepository = (*MongoDrafts)(nil)

type MongoDrafts struct {
	collection *mongo.Collection
}

func NewMongoDrafts(db *mongo.Database) *MongoDrafts {
	return &MongoDrafts{db.Collection("drafts")}
}

func (r *MongoDrafts) Update(ctx context.Context, draft *Draft) error {
	_, err := r.collection.UpdateOne(ctx, map[string]any{"_id": draft.ID}, map[string]any{"$set": draft})
	return err
}

func (r *MongoDrafts) GetByConfigID(ctx context.Context, configID primitive.ObjectID) (*Draft, error) {
	var draft Draft
	err := r.collection.FindOne(ctx, map[string]any{"configId": configID}).Decode(&draft)
	switch err {
	case nil:
		return &draft, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}
