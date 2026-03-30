package routes

import (
	"fmt"

	"github.com/yaien/cultural/internal/infrastructure"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/controllers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func auth(mono *infrastructure.Monolith, app *application.Application) {

	config := &oauth2.Config{
		ClientID:     mono.Config.Google.ClientID,
		ClientSecret: mono.Config.Google.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", mono.Config.Server.URL),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     endpoints.Google,
	}

	ctrl := controllers.NewAuthController(app.Deps.Auth.Accounts, mono.SessionStore, config)
	mono.WebRouter.HandleFunc("GET /auth/google/login", ctrl.Login)
	mono.WebRouter.HandleFunc("GET /auth/google/callback", ctrl.Callback)
	mono.WebRouter.HandleFunc("POST /auth/logout", ctrl.Logout)
}
