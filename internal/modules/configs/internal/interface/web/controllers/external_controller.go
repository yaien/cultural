package controllers

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
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
		http.Error(w, "invalid organization id", http.StatusBadRequest)
		return
	}

	filename := r.PathValue("filename")

	file, data, err := c.app.GetFile(r.Context(), organizationID, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	w.WriteHeader(http.StatusOK)

	_, err = bufio.NewWriter(w).ReadFrom(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
