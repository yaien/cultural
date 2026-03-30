package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	repository Repository
}

func NewUsers(repository Repository) *Users {
	return &Users{repository: repository}
}

func (c *Users) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	return c.repository.GetByID(ctx, id)
}
