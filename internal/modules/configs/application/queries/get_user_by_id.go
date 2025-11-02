package queries

import (
	"context"

	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserByIDQuery struct {
	users models.UserRepository
}

func NewGetUserByIDQuery(users models.UserRepository) *GetUserByIDQuery {
	return &GetUserByIDQuery{users: users}
}

func (q *GetUserByIDQuery) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return q.users.GetByID(ctx, id)
}
