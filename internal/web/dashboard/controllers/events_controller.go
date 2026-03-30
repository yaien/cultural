package controllers

import (
	"net/http"

	"github.com/yaien/cultural/internal/web/dashboard/views/events"
)

type EventsController struct {
}

func NewEventsController() *EventsController {
	return &EventsController{}
}

func (c *EventsController) Index(w http.ResponseWriter, r *http.Request) {
	_ = events.Page().Render(r.Context(), w)
}
