package store

import (
	"context"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/coderror"
)

var _ interface {
	integration.Definition
	integration.Template
} = (*Integration)(nil)

type Integration struct {
	store *Store
}

func NewIntegration(s *Store) *Integration {
	return &Integration{store: s}
}

func (s *Integration) Name() string { return "store" }

func (s *Integration) TemplateFuncMap(ctx context.Context, config *label.Config) integration.FuncMap {
	return integration.FuncMap{}
}

func (s *Integration) TemplatePresetMap(ctx context.Context, config *label.Config) integration.PresetMap {
	return integration.PresetMap{
		"product": {
			Key:  "product",
			Name: "Producto",
			Options: func() (options []label.Option, err error) {
				products, err := s.store.Products.GetByOrganizationID(ctx, config.OrganizationID)
				if err != nil {
					return nil, err
				}
				for _, product := range products {
					options = append(options, label.Option{
						Value: product.Slug,
						Label: product.Name,
					})
				}
				return options, nil
			},
			Load: func(params ...string) (any, error) {
				if len(params) == 0 {
					return nil, coderror.Newf("missing_param", "misssing param product id")
				}

				product, err := s.store.Products.GetBySlugAndOrganizationID(ctx, params[0], config.OrganizationID)
				if err != nil {
					return nil, coderror.Newf("product_not_found", "product not found: %w", err)
				}

				preset := struct {
					Product *Product
				}{&product}

				return preset, nil
			},
		},
	}
}
