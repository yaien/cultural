package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvitationController struct {
	app *application.Application
}

func NewInvitationController(app *application.Application) *InvitationController {
	return &InvitationController{app: app}
}

func (c *InvitationController) Accept(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)
	user := r.Context().Value(middlewares.UserContextKey).(*models.User)

	err = c.app.AcceptInvitation(r.Context(), &commands.AcceptInvitationRequest{
		InvitationID:   id,
		OrganizationID: config.OrganizationID,
		UserID:         user.ID,
		UserEmail:      user.Email,
		UserName:       user.Name,
		UserAvatarUrl:  user.AvatarUrl,
	})

	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed accepting invitation: %w", err))
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)

}
