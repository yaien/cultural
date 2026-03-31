package auth

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
)

type Account struct {
	Provider     string
	ID           string
	Name         string
	AvatarUrl    string
	Email        string
	AcccessToken string
	RefreshToken string
	ExpiresAt    time.Time
	LastUsedAt   time.Time
}

type User struct {
	ID        primitive.ID `gorm:"primaryKey;autoIncrement"`
	Email     string       `gorm:"uniqueIndex"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	AvatarUrl string
	Accounts  map[string]*Account `gorm:"type:jsonb;serializer:json"`
}

type Repository interface {
	GetByID(ctx context.Context, id primitive.ID) (*User, error)
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
