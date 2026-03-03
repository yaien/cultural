package models

import (
	"context"
	"time"

	"github.com/markbates/goth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Email     string               `bson:"email" json:"email"`
	Name      string               `bson:"name" json:"name"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
	AvatarUrl string               `bson:"avatarUrl" json:"avatarUrl"`
	Accounts  map[string]goth.User `bson:"accounts" json:"accounts"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
