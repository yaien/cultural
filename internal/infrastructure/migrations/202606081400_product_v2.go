package migrations

import (
	"context"

	"gorm.io/gorm"
)

func init() {
	Register(Migration{
		Name: "202606081400_product_v2",
		Up: func(ctx context.Context, db *gorm.DB) error {
			type Product struct {
				DeletedAt gorm.DeletedAt `gorm:"index"`
			}

			return db.AutoMigrate(&Product{})
		},
	})
}
