package instagram

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

//go:generate go tool templ generate

func (i *Instagram) Page(ctx context.Context, config *models.Config) (templ.Component, error) {
	integration, err := i.integrations.GetByOrganizationIDAndName(ctx, config.OrganizationID, i.Name())
	if err != nil && !models.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed at get integration: %w", err)
	}

	return Page(integration), nil
}
