package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/home"
)

type DashboardController struct {
	app *application.Application
}

func NewDashboardController(app *application.Application) *DashboardController {
	return &DashboardController{app: app}
}

func (c *DashboardController) Home(w http.ResponseWriter, r *http.Request) {
	_ = home.Home().Render(r.Context(), w)
}

func (c *DashboardController) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
