package label

import (
	"context"
	"time"

	"github.com/yaien/cultural/internal/lib/primitive"
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

type ConfigRepository interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
	GetByOrganizationID(ctx context.Context, id primitive.ID) (*Config, error)
	Update(ctx context.Context, config *Config) error
}

type Configs struct {
	configs ConfigRepository
}

func NewConfigs(configs ConfigRepository) *Configs {
	return &Configs{configs: configs}
}

func (c *Configs) GetByHost(ctx context.Context, host string) (*Config, error) {
	return c.configs.GetByHost(ctx, host)
}

func (c *Configs) GetByOrganizationID(ctx context.Context, organizationID primitive.ID) (*Config, error) {
	return c.configs.GetByOrganizationID(ctx, organizationID)
}

func (c *Configs) Update(ctx context.Context, config *Config) error {
	return c.configs.Update(ctx, config)
}
