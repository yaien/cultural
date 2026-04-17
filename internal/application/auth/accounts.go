package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Accounts struct {
	repository gorm.Interface[User]
}

func NewAccounts(db *gorm.DB) *Accounts {
	return &Accounts{gorm.G[User](db)}
}

func (c *Accounts) Sync(ctx context.Context, account *Account) (User, error) {
	user, err := c.repository.Where("email = ?", account.Email).Take(ctx)

	switch {
	case err == nil:

		user.Accounts[account.Provider] = account

		if _, err = c.repository.Updates(ctx, user); err != nil {
			return user, fmt.Errorf("failed updating user: %w", err)
		}

		return user, nil

	case errors.Is(err, gorm.ErrRecordNotFound):

		user = User{
			Email:     account.Email,
			Name:      account.Name,
			AvatarUrl: account.AvatarUrl,
			Accounts:  map[string]*Account{account.Provider: account},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = c.repository.Create(ctx, &user)
		if err != nil {
			return user, fmt.Errorf("failed creating user: %w", err)
		}

		return user, nil

	default:
		return user, fmt.Errorf("failed getting user by email: %w", err)
	}
}
