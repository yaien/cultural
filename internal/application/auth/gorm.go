package auth

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ Repository = (*Gorm)(nil)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{db: db}
}

func (r *Gorm) GetByID(ctx context.Context, id primitive.ID) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, primitive.Error(err)
}

func (r *Gorm) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, primitive.Error(err)
}

func (r *Gorm) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Gorm) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
