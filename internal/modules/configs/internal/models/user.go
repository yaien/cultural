package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	Provider     string    `bson:"provider" json:"provider"`
	ID           string    `bson:"id" json:"id"`
	Name         string    `bson:"name" json:"name"`
	AvatarUrl    string    `bson:"avatarUrl" json:"avatarUrl"`
	Email        string    `bson:"email" json:"email"`
	AcccessToken string    `bson:"accessToken" json:"accessToken"`
	RefreshToken string    `bson:"refreshToken" json:"refreshToken"`
	ExpiresAt    time.Time `bson:"expiresAt" json:"expiresAt"`
	LastUsedAt   time.Time `bson:"lastUsedAt" json:"lastUsedAt"`
}

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Email     string              `bson:"email" json:"email"`
	Name      string              `bson:"name" json:"name"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`
	AvatarUrl string              `bson:"avatarUrl" json:"avatarUrl"`
	Accounts  map[string]*Account `bson:"accounts" json:"accounts"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
