package repositories

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ models.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	switch err {
	case nil:
		return &user, nil
	case mongo.ErrNoDocuments:
		return nil, models.NotFoundError(err)
	default:
		return nil, err
	}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	switch err {
	case nil:
		return &user, nil
	case mongo.ErrNoDocuments:
		return nil, models.NotFoundError(err)
	default:
		return nil, err
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}
