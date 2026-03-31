package integration

import (
	"context"

	"github.com/yaien/cultural/internal/lib/primitive"

	"gorm.io/gorm"
)

var _ Repository[any] = (*Gorm[any])(nil)

type Gorm[T any] struct {
	db *gorm.DB
}

func NewGorm[T any](db *gorm.DB) *Gorm[T] {
	return &Gorm[T]{db: db}
}

func (r *Gorm[T]) Create(ctx context.Context, i *Integration[T]) error {
	return r.db.WithContext(ctx).Create(i).Error
}

func (r *Gorm[T]) Update(ctx context.Context, i *Integration[T]) error {
	return r.db.WithContext(ctx).Save(i).Error
}

func (r *Gorm[T]) GetByOrganizationIDAndName(ctx context.Context, organizationID primitive.ID, name string) (*Integration[T], error) {
	var i Integration[T]
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND name = ?", organizationID, name).
		First(&i).Error
	return &i, primitive.Error(err)
}

func (r *Gorm[T]) GetByName(ctx context.Context, name string) ([]*Integration[T], error) {
	var integrations []*Integration[T]
	err := r.db.WithContext(ctx).Where("name = ?", name).Find(&integrations).Error
	if err != nil {
		return nil, err
	}
	return integrations, nil
}
