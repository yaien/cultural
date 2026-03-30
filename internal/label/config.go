package label

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	OrganizationID primitive.ObjectID `bson:"organizationId"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
	Host           string             `bson:"host"`
	Title          string             `bson:"title"`
	Url            string             `bson:"url"`
	Email          string             `bson:"email"`
	Fonts          map[string]*Font   `bson:"fonts"`
	Pages          map[string]*Page   `bson:"pages"`
	Layouts        map[string]*Layout `bson:"layouts"`
	Emails         map[string]*Email  `bson:"emails"`
	Colors         []*Color           `bson:"colors"`
}

type ConfigRepository interface {
	GetByHost(ctx context.Context, host string) (*Config, error)
	GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*Config, error)
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

func (c *Configs) GetByOrganizationID(ctx context.Context, organizationID primitive.ObjectID) (*Config, error) {
	return c.configs.GetByOrganizationID(ctx, organizationID)
}

func (c *Configs) Update(ctx context.Context, config *Config) error {
	return c.configs.Update(ctx, config)
}
