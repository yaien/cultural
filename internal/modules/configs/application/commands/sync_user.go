package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/markbates/goth"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type SyncUserCommand struct {
	users models.UserRepository
}

func NewSyncUserCommand(users models.UserRepository) *SyncUserCommand {
	return &SyncUserCommand{users: users}
}

func (c *SyncUserCommand) SyncUser(ctx context.Context, user goth.User) (*models.User, error) {
	u, err := c.users.GetByEmail(ctx, user.Email)
	var e *models.Error
	switch {
	case err == nil:

		u.Accounts[user.Provider] = user

		err = c.users.Update(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed updating user: %w", err)
		}

		return u, nil

	case errors.As(err, &e) && e.Code == "not_found":

		u = &models.User{
			Email:     user.Email,
			Name:      user.Name,
			AvatarUrl: user.AvatarURL,
			Accounts:  map[string]goth.User{user.Provider: user},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = c.users.Create(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed creating user: %w", err)
		}

		return u, nil

	default:
		return nil, fmt.Errorf("failed getting user by email: %w", err)
	}
}
