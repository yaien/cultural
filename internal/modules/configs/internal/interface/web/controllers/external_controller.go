package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/library/storage"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalController struct {
	storage *storage.Storage
}

func NewExternalController(s *storage.Storage) *ExternalController {
	return &ExternalController{s}
}

func (c *ExternalController) GetFile(w http.ResponseWriter, r *http.Request) {
	var req storage.DownloadOptions
	var err error

	req.OrganizationID, err = primitive.ObjectIDFromHex(r.PathValue("organization_id"))
	if err != nil {
		WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid organization id: %w", err)))
		return
	}

	req.Name = r.PathValue("filename")

	if variant := r.URL.Query().Get("variant"); variant != "" {
		if req.Variant, err = strconv.Atoi(variant); err != nil {
			WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid quality: %w", err)))
			return
		}
	}

	ctx := r.Context()

	res, err := c.storage.Download(ctx, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteFile(w, r, res)
}
