package instagram

import (
	"context"
	"text/template"

	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

func (i *Instagram) TemplateFuncMap(ctx context.Context, config *models.Config) template.FuncMap {
	return template.FuncMap{
		"get_instagram_posts": func() ([]*Post, error) {
			integration, err := i.integrations.GetByOrganizationIDAndName(ctx, config.OrganizationID, i.Name())
			if err != nil {
				return nil, err
			}

			if integration == nil || !integration.Data.Connected {
				return nil, nil
			}

			return integration.Data.Posts, nil
		},
	}

}
