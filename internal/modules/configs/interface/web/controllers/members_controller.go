package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
)

type MembersController struct {
	app *application.Application
}

func NewMembersController(app *application.Application) *MembersController {
	return &MembersController{app: app}
}

func (c *MembersController) Index(w http.ResponseWriter, r *http.Request) {
	views.Members(w, r)
}
