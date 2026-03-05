package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views"
)

type DashboardController struct {
	app *application.Application
}

func NewDashboardController(app *application.Application) *DashboardController {
	return &DashboardController{app: app}
}

func (c *DashboardController) Home(w http.ResponseWriter, r *http.Request) {
	_ = views.Home().Render(r.Context(), w)
}

func (c *DashboardController) Toast(w http.ResponseWriter, r *http.Request) {
	toast, exists, err := GetToast(w, r)
	if err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed to get toast: %w", err))
		return
	}

	if !exists {
		w.WriteHeader(http.StatusOK)
		return
	}

	_ = views.Toast(toast.Message, toast.Variant).Render(r.Context(), w)
}

func (c *DashboardController) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
