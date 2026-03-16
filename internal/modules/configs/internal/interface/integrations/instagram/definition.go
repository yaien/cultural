package instagram

import (
	"time"

	"github.com/yaien/cultural/internal/modules/configs/internal/interface/repositories"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ interface {
	models.IntegrationDefinition
	models.IntegrationOAuth
	models.IntegrationTemplateFuncMap
	models.IntegrationBackground
} = (*Instagram)(nil)

type Data struct {
	Connected bool      `bson:"connected"`
	User      *User     `bson:"user"`
	Posts     []*Post   `bson:"posts"`
	Token     string    `bson:"token"`
	ExpireAt  time.Time `bson:"expireAt"`
}

type Instagram struct {
	integrations models.IntegrationRepository[Data]
	configs      models.ConfigRepository
}

func Mew(db *mongo.Database) *Instagram {
	return &Instagram{
		integrations: repositories.NewIntegrationRepository[Data](db),
		configs:      repositories.NewConfigRepository(db),
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
