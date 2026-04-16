package instagram

import (
	"context"
	"errors"
	"fmt"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/application/label"
	"gorm.io/gorm"
)

//go:generate go tool templ generate

func (i *Instagram) Page(ctx context.Context, config *label.Config) (templ.Component, error) {
	integration, err := i.integrations.Where("organization_id = ? AND name = ?", config.OrganizationID, i.Name()).First(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed at get integration: %w", err)
	}

	return Page(&integration), nil
}
