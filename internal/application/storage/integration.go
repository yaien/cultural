package storage

import (
	"context"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
)

var _ interface {
	integration.Definition
	integration.TemplateFuncMapper
} = (*Integration)(nil)

type Integration struct {
}

func NewIntegration() *Integration {
	return &Integration{}
}

func (i *Integration) Name() string {
	return "storage"
}

func (i *Integration) TemplateFuncMap(ctx context.Context, config *label.Config) integration.TemplateFuncMap {
	return integration.TemplateFuncMap{
		"file_url":          FileURL,
		"external_file_url": NewExternalURLFunc(config.Url, config.OrganizationID),
	}
}
