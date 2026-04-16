package auth

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

type Users struct {
	repository gorm.Interface[User]
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{gorm.G[User](db)}
}

func (c *Users) GetByID(ctx context.Context, id primitive.ID) (User, error) {
	return c.repository.Where("id = ?", id).Take(ctx)
}
