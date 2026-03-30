package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	Provider     string    `bson:"provider"`
	ID           string    `bson:"id"`
	Name         string    `bson:"name"`
	AvatarUrl    string    `bson:"avatarUrl"`
	Email        string    `bson:"email"`
	AcccessToken string    `bson:"accessToken"`
	RefreshToken string    `bson:"refreshToken"`
	ExpiresAt    time.Time `bson:"expiresAt"`
	LastUsedAt   time.Time `bson:"lastUsedAt"`
}

type User struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Email     string              `bson:"email"`
	Name      string              `bson:"name"`
	CreatedAt time.Time           `bson:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt"`
	AvatarUrl string              `bson:"avatarUrl"`
	Accounts  map[string]*Account `bson:"accounts"`
}

type Repository interface {
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type Auth struct {
	Users    *Users
	Accounts *Accounts
}

func New(repo Repository) *Auth {
	return &Auth{
		Users:    NewUsers(repo),
		Accounts: NewAccounts(repo),
	}
}
