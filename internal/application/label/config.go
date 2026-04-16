package label

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
	"gorm.io/gorm"
)

type Config struct {
	ID             primitive.ID `gorm:"primaryKey,autoIncrement"`
	OrganizationID primitive.ID `gorm:"index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Host           string `gorm:"index"`
	Title          string
	Url            string
	Email          string
	Fonts          map[string]*Font   `gorm:"type:jsonb;serializer:json"`
	Pages          map[string]*Page   `gorm:"type:jsonb;serializer:json"`
	Layouts        map[string]*Layout `gorm:"type:jsonb;serializer:json"`
	Emails         map[string]*Email  `gorm:"type:jsonb;serializer:json"`
	Colors         []*Color           `gorm:"type:jsonb;serializer:json"`
}

type Configs struct {
	configs gorm.Interface[Config]
}

func NewConfigs(db *gorm.DB) *Configs {
	return &Configs{gorm.G[Config](db)}
}

func (c *Configs) GetByHost(ctx context.Context, host string) (Config, error) {
	return c.configs.Where("host = ?", host).Take(ctx)
}

func (c *Configs) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) (Config, error) {
	return c.configs.Where("organization_id = ?", organizationID).Take(ctx)
}

func (c *Configs) Update(ctx context.Context, config Config) error {
	_, err := c.configs.Updates(ctx, config)
	return err
}
