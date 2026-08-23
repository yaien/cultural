package instagram

import (
	"context"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
)

func (i *Instagram) TemplateFuncMap(ctx context.Context, config *label.Config) integration.TemplateFuncMap {
	return integration.TemplateFuncMap{
		"get_instagram_posts": func() ([]*Post, error) {
			integration, err := i.integrations.
				Where("organization_id = ? and name = ?", config.OrganizationID, i.Name()).
				First(ctx)

			if err != nil {
				return nil, err
			}

			if !integration.Data.Connected {
				return nil, nil
			}

			return integration.Data.Posts, nil
		},
	}
}
