package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/web/dashboard/views/home"
)

type DashboardController struct {
}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (c *DashboardController) Home(w http.ResponseWriter, r *http.Request) {
	_ = home.Home().Render(r.Context(), w)
}

func (c *DashboardController) Empty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
