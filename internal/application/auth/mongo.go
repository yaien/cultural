package auth

import (
	"context"

	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ Repository = (*Mongo)(nil)

type Mongo struct {
	db *mongo.Database
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (r *Mongo) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	switch err {
	case nil:
		return &user, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *Mongo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	switch err {
	case nil:
		return &user, nil
	case mongo.ErrNoDocuments:
		return nil, coderror.New(coderror.NotFound, err)
	default:
		return nil, err
	}
}

func (r *Mongo) Create(ctx context.Context, user *User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	return err
}

func (r *Mongo) Update(ctx context.Context, user *User) error {
	_, err := r.db.Collection("users").UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}
