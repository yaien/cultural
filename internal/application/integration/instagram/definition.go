package instagram

import (
	"time"

	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ interface {
	integration.Definition
	integration.TemplateFuncMap
	integration.OAuth
	integration.Background
} = (*Instagram)(nil)

type Data struct {
	Connected bool      `bson:"connected"`
	User      *User     `bson:"user"`
	Posts     []*Post   `bson:"posts"`
	Token     string    `bson:"token"`
	ExpireAt  time.Time `bson:"expireAt"`
}

type Instagram struct {
	integrations integration.Repository[Data]
	configs      label.ConfigRepository
}

func Mew(db *mongo.Database) *Instagram {
	return &Instagram{
		integrations: integration.NewMongo[Data](db),
		configs:      label.NewMongoConfigs(db),
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
