package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/views"
)

type DashboardController struct {
	app *application.Application
}

func NewDashboardController(app *application.Application) *DashboardController {
	return &DashboardController{app: app}
}

func (c *DashboardController) Home(w http.ResponseWriter, r *http.Request) {
	views.Home(w, r)
}
