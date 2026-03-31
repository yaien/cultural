package store

import (
	"context"
	"errors"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

var _ Repository = (*Gorm)(nil)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{db: db}
}

func (r *Gorm) Create(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *Gorm) Update(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *Gorm) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) (products []*Product, err error) {
	err = r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&products).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if products == nil {
		products = []*Product{}
	}

	return products, nil
}

func (r *Gorm) GetByIDAndOrganizationID(ctx context.Context, productID, organizationID primitive.ID) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).Where("id = ? AND organization_id = ?", productID, organizationID).First(&product).Error
	return &product, primitive.Error(err)
}

func (r *Gorm) GetBySlugAndOrganizationID(ctx context.Context, slug string, organizationID primitive.ID) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).Where("slug = ? AND organization_id = ?", slug, organizationID).First(&product).Error
	return &product, primitive.Error(err)
}
