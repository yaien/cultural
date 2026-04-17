package worker

import (
	"context"

	"gorm.io/gorm"
)

type Gorm struct {
	db    *gorm.DB
	limit int64
}

var _ Store = (*Gorm)(nil)

// NewGormStore creates a new Gorm store using Job as the model.
func NewGormStore(db *gorm.DB) *Gorm {
	return &Gorm{
		db:    db,
		limit: 10,
	}
}

func (g *Gorm) Create(ctx context.Context, job Job) error {
	return g.db.WithContext(ctx).Create(&job).Error
}

func (g *Gorm) Update(ctx context.Context, job Job) error {
	return g.db.WithContext(ctx).Save(&job).Error
}

func (g *Gorm) Fetch(ctx context.Context) (jobs []Job, err error) {
	err = g.db.WithContext(ctx).
		Where("status = ?", StatusPending).
		Order("created_at ASC").
		Limit(int(g.limit)).
		Find(&jobs).Error

	return
}
