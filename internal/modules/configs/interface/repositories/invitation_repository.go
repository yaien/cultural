package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.InvitationRepository = (*InvitationRepository)(nil)

type InvitationRepository struct {
	db *mongo.Database
}

func NewInvitationRepository(db *mongo.Database) *InvitationRepository {
	return &InvitationRepository{
		db: db,
	}
}

func (r *InvitationRepository) Create(ctx context.Context, invitation *models.Invitation) error {

	res, err := r.db.Collection("invitations").InsertOne(ctx, invitation)
	if err != nil {
		return err
	}

	invitation.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}
