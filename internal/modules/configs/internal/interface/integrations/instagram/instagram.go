package instagram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/a-h/templ"
	"github.com/spf13/viper"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/repositories"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var _ interface {
	models.IntegrationDefinition
	models.IntegrationOAuth
} = (*Instagram)(nil)

type Data struct {
}

type Instagram struct {
	integrations models.IntegrationRepository[Data]
}

func Mew(db *mongo.Database) *Instagram {
	return &Instagram{
		integrations: repositories.NewIntegrationRepository[Data](db),
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

func (i *Instagram) Page(ctx context.Context, config *models.Config) (templ.Component, error) {
	integration, err := i.integrations.Get(ctx, models.GetIntegrationOptions{
		OrganizationID: config.OrganizationID,
		Name:           i.Name(),
	})

	if err != nil && !models.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed at get integration: %w", err)
	}

	return Page(integration), nil
}

func (i *Instagram) OAuthConfig(config *models.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("INSTAGRAM_CLIENT_ID"),
		ClientSecret: viper.GetString("INSTAGRAM_CLIENT_SECRET"),
		Endpoint:     endpoints.Instagram,
		Scopes:       []string{"instagram_business_basic"},
		RedirectURL:  fmt.Sprintf("%s/dashboard/integrations/instagram/oauth/callback", config.Url),
	}
}

func (i *Instagram) OAuthCodeURL(ctx context.Context, config *models.Config) (string, error) {
	return i.OAuthConfig(config).AuthCodeURL(""), nil
}

func (i *Instagram) OAuthExchange(ctx context.Context, config *models.Config, code string) error {
	token, err := i.OAuthConfig(config).Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed getting token:%w", err)
	}

	slog.Info("got insta access token",
		"token", token.AccessToken,
		"expires in", (time.Duration(token.ExpiresIn) * time.Second).String(),
		"expiration", token.Expiry.Format(time.RFC3339),
		"refresh token", token.RefreshToken,
	)

	return nil

}
