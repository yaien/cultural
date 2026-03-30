package instagram

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/coderror"
)

//go:generate go tool templ generate

func (i *Instagram) Page(ctx context.Context, config *label.Config) (templ.Component, error) {
	integration, err := i.integrations.GetByOrganizationIDAndName(ctx, config.OrganizationID, i.Name())
	if err != nil && !coderror.Is(err, coderror.NotFound) {
		return nil, fmt.Errorf("failed at get integration: %w", err)
	}

	return Page(integration), nil
}
