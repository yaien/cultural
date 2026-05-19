package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/application/label"
	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"

	"github.com/yaien/cultural/internal/web/dashboard/views/pages"
	"github.com/yaien/cultural/internal/web/middlewares"
)

type FilesController struct {
	storage *storage.Storage
}

func NewFilesController(s *storage.Storage) *FilesController {
	return &FilesController{s}
}

func (fc *FilesController) Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(5 << 20) // 5 MB
	if err != nil {
		WriteJSONErr(w, coderror.Newf(coderror.DecodeFailed, "failed to parse multipart form: %w", err))
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)

	var files []storage.File

	for _, handler := range r.MultipartForm.File["files"] {
		data, err := handler.Open()
		if err != nil {
			WriteHTMLErr(w, fmt.Errorf("failed opening file: %w", err))
			return
		}

		defer func() {
			if err := data.Close(); err != nil {
				slog.Warn("failed closing file data", "err", err)
			}
		}()

		file, err := fc.storage.Upload(ctx, &storage.UploadOptions{
			Name:           handler.Filename,
			Size:           handler.Size,
			ContentType:    handler.Header.Get("Content-Type"),
			OrganizationID: config.OrganizationID,
			Data:           data,
		})

		files = append(files, file)

		if err != nil {
			WriteHTMLErr(w, fmt.Errorf("failed uploading file: %w", err))
		}
	}
	w.Header().Set("HX-Trigger", "render")
	_ = pages.FileGrid(files, storage.FileURL).Render(ctx, w)

}

func (fc *FilesController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
	filename := r.PathValue("filename")

	file, err := fc.storage.GetByOrganizationIDAndName(ctx, config.OrganizationID, filename)
	if err != nil {
		WriteHTMLErr(w, err)
		return
	}

	if err := fc.storage.Delete(ctx, config.OrganizationID, file.ID); err != nil {
		WriteHTMLErr(w, err)
		return
	}

	w.Header().Set("HX-Trigger", "deleted, render")
	w.WriteHeader(http.StatusOK)
}

func (fc *FilesController) Download(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*label.Config)

	var err error
	var req storage.DownloadOptions

	req.Name = r.PathValue("filename")
	req.OrganizationID = config.OrganizationID

	if id, err := primitive.ParseID(req.Name); err == nil {
		req.ID = &id
		req.Name = ""
	}

	if variant := r.URL.Query().Get("variant"); variant != "" {
		if req.Variant, err = strconv.Atoi(variant); err != nil {
			WriteJSONErr(w, coderror.Newf(coderror.DecodeFailed, "invalid variant: %w", err))
			return
		}
	}

	res, err := fc.storage.Download(r.Context(), &req)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed getting file: %w", err))
		return
	}

	WriteFile(w, r, res)

}

func (fc *FilesController) Rename(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*label.Config)
	filename := r.PathValue("filename")

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, coderror.Newf(coderror.DecodeFailed, "failed parsing form: %w", err))
		return
	}

	newName := r.PostFormValue("name")

	if err := fc.storage.Rename(ctx, config.OrganizationID, filename, newName); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed renaming file: %w", err))
		return
	}

	w.Header().Set("HX-Trigger", "renamed, render")
	w.WriteHeader(http.StatusOK)
}
