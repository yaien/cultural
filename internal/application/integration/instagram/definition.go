package instagram

import (
	"time"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"gorm.io/gorm"
)

var _ interface {
	integration.Definition
	integration.TemplateFuncMap
	integration.OAuth
	integration.Background
} = (*Instagram)(nil)

type Data struct {
	Connected bool
	User      *User
	Posts     []*Post
	Token     string
	ExpireAt  time.Time
}

type Integration = integration.Integration[Data]

type Instagram struct {
	integrations gorm.Interface[Integration]
	configs      *label.Configs
}

func New(db *gorm.DB, configs *label.Configs) *Instagram {
	return &Instagram{
		integrations: gorm.G[Integration](db),
		configs:      configs,
	}
}

func (i *Instagram) Title() string {
	return "Instagram"
}

func (i *Instagram) Name() string {
	return "instagram"
}

func (i *Instagram) Description() string {
	return "Trae tus posts de instagram a tu web"
}

func (i *Instagram) Image() string {
	return "instagram.png"
}
