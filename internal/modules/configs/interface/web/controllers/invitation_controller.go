package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvitationController struct {
	App *application.Application
}

func NewInvitationController(app *application.Application) *InvitationController {
	return &InvitationController{App: app}
}

func (c *InvitationController) OnInvitation(w http.ResponseWriter, r *http.Request) {
	_, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_ = views.OnInvitation().Render(r.Context(), w)
}
