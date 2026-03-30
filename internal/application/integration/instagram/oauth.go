package instagram

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/yaien/cultural/internal/application/integration"
	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/lib/coderror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func (i *Instagram) OAuthConfig(config *label.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("INSTAGRAM_CLIENT_ID"),
		ClientSecret: viper.GetString("INSTAGRAM_CLIENT_SECRET"),
		Endpoint:     endpoints.Instagram,
		Scopes:       []string{"instagram_business_basic"},
		RedirectURL:  fmt.Sprintf("%s/dashboard/integrations/instagram/oauth/callback", config.Url),
	}
}

func (i *Instagram) OAuthCodeURL(ctx context.Context, config *label.Config) (string, error) {
	return i.OAuthConfig(config).AuthCodeURL(""), nil
}

func (i *Instagram) OAuthExchange(ctx context.Context, config *label.Config, code string) error {
	auth := i.OAuthConfig(config)

	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed getting token:%w", err)
	}

	client := NewClient(token, auth)

	token, err = client.GetLongToken(ctx)
	if err != nil {
		return fmt.Errorf("failed getting long token")
	}

	client.SetToken(token)

	user, err := client.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("failed getting user: %w", err)
	}

	posts, err := client.GetPosts(ctx)
	if err != nil {
		return fmt.Errorf("failed getting posts: %w", err)
	}

	if len(posts) > 10 {
		posts = posts[:6]
	}

	data := Data{
		Connected: true,
		User:      user,
		Posts:     posts,
		Token:     token.AccessToken,
		ExpireAt:  time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}

	if err := i.Save(ctx, config.OrganizationID, data); err != nil {
		return fmt.Errorf("failed saving data: %w", err)
	}

	return nil

}

func (i *Instagram) Save(ctx context.Context, organizationID primitive.ObjectID, data Data) error {
	itg, err := i.integrations.GetByOrganizationIDAndName(ctx, organizationID, i.Name())

	switch {
	case err == nil:
		itg.Data = data
		itg.UpdatedAt = time.Now()
		return i.integrations.Update(ctx, itg)

	case coderror.Is(err, coderror.NotFound):
		itg = &integration.Integration[Data]{
			ID:             primitive.NewObjectID(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Name:           i.Name(),
			OrganizationID: organizationID,
			Data:           data,
		}

		return i.integrations.Create(ctx, itg)

	default:
		return fmt.Errorf("failed getting integration: %w", err)
	}
}
