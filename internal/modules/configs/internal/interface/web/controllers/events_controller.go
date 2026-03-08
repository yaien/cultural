package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/events"
)

type EventsController struct {
	app *application.Application
}

func NewEventsController(app *application.Application) *EventsController {
	return &EventsController{
		app: app,
	}
}

func (c *EventsController) Index(w http.ResponseWriter, r *http.Request) {
	_ = events.Page().Render(r.Context(), w)
}
