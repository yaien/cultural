package instagram

import (
	"time"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"gorm.io/gorm"
)

var _ interface {
	integration.Definition
	integration.TemplateFuncMapper
	integration.OAuth
	integration.Background
	integration.Page
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

func (i *Instagram) Name() string {
	return "instagram"
}
