package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/yaien/cultural/internal/lib/coderror"
)

type Accounts struct {
	repository Repository
}

func NewAccounts(repository Repository) *Accounts {
	return &Accounts{repository: repository}
}

func (c *Accounts) Sync(ctx context.Context, account *Account) (*User, error) {
	u, err := c.repository.GetByEmail(ctx, account.Email)

	switch {
	case err == nil:

		u.Accounts[account.Provider] = account

		err = c.repository.Update(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed updating user: %w", err)
		}

		return u, nil

	case coderror.Is(err, coderror.NotFound):

		u = &User{
			Email:     account.Email,
			Name:      account.Name,
			AvatarUrl: account.AvatarUrl,
			Accounts:  map[string]*Account{account.Provider: account},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = c.repository.Create(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed creating user: %w", err)
		}

		return u, nil

	default:
		return nil, fmt.Errorf("failed getting user by email: %w", err)
	}
}
