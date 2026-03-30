package admin

import (
	"context"

	"github.com/yaien/cultural/internal/library/coderror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ InvitationRepository = (*MongoInvitations)(nil)

type MongoInvitations struct {
	collection *mongo.Collection
}

func NewMongoInvitations(db *mongo.Database) *MongoInvitations {
	return &MongoInvitations{db.Collection("invitations")}
}

func (r *MongoInvitations) GetByIDAndOrganizationID(ctx context.Context, id, organizationID primitive.ObjectID) (*Invitation, error) {
	var invitation Invitation
	err := r.collection.FindOne(ctx, primitive.M{"_id": id, "organizationId": organizationID}).Decode(&invitation)
	switch err {
	case nil:
		return &invitation, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *MongoInvitations) Create(ctx context.Context, invitation *Invitation) error {
	_, err := r.collection.InsertOne(ctx, invitation)
	return err
}

func (r *MongoInvitations) Update(ctx context.Context, invitation *Invitation) error {
	_, err := r.collection.UpdateOne(ctx, primitive.M{"_id": invitation.ID}, primitive.M{"$set": invitation})
	return err
}
