package auth

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Users struct {
	repository Repository
}

func NewUsers(repository Repository) *Users {
	return &Users{repository: repository}
}

func (c *Users) GetByID(ctx context.Context, id primitive.ID) (*User, error) {
	return c.repository.GetByID(ctx, id)
}
