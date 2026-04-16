package auth

import (
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
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

type Auth struct {
	Users    *Users
	Accounts *Accounts
}

func New(db *gorm.DB) *Auth {
	return &Auth{
		Users:    NewUsers(db),
		Accounts: NewAccounts(db),
	}
}
