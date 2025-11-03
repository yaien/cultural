package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvitationController struct {
	App *application.Application
}

func NewInvitationController(app *application.Application) *InvitationController {
	return &InvitationController{App: app}
}

func (c *InvitationController) OnInvitation(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	user := r.Context().Value(models.UserContextKey).(*models.User)

	err = c.App.AcceptInvitation(r.Context(), &commands.AcceptInvitationRequest{
		InvitationID:   id,
		OrganizationID: config.OrganizationID,
		UserID:         user.ID,
		UserEmail:      user.Email,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = views.OnInvitation().Render(r.Context(), w)

}
