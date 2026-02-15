package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
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

func (r *InvitationRepository) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*models.Invitation, error) {
	var invitation models.Invitation
	err := r.db.Collection("invitations").FindOne(ctx, primitive.M{"_id": id, "organizationId": organizationID}).Decode(&invitation)
	return &invitation, translate(err)
}

func (r *InvitationRepository) Create(ctx context.Context, invitation *models.Invitation) error {
	res, err := r.db.Collection("invitations").InsertOne(ctx, invitation)
	if err != nil {
		return err
	}

	invitation.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *InvitationRepository) Update(ctx context.Context, invitation *models.Invitation) error {
	_, err := r.db.Collection("invitations").UpdateOne(ctx, primitive.M{"_id": invitation.ID}, primitive.M{"$set": invitation})
	return err
}
