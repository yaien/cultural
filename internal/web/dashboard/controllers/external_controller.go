package controllers

import (
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
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

	req.OrganizationID, err = primitive.ParseID(r.PathValue("organization_id"))
	if err != nil {
		WriteJSONErr(w, coderror.Newf(coderror.DecodeFailed, "invalid organization id: %w", err))
		return
	}

	req.Name = r.PathValue("filename")

	if variant := r.URL.Query().Get("variant"); variant != "" {
		if req.Variant, err = strconv.Atoi(variant); err != nil {
			WriteJSONErr(w, coderror.Newf(coderror.DecodeFailed, "invalid variant: %w", err))
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
