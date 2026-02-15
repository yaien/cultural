package controllers

import (
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalController struct {
	app *application.Application
}

func NewExternalController(app *application.Application) *ExternalController {
	return &ExternalController{
		app: app,
	}
}

func (c *ExternalController) GetFile(w http.ResponseWriter, r *http.Request) {
	organizationID, err := primitive.ObjectIDFromHex(r.PathValue("organization_id"))
	if err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid organization id: %w", err)))
		return
	}

	filename := r.PathValue("filename")

	file, data, err := c.app.GetFile(r.Context(), organizationID, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteFile(w, file, data)
}
