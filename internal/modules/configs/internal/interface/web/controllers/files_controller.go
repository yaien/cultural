package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/pages"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type FilesController struct {
	app *application.Application
}

func NewFilesController(app *application.Application) *FilesController {
	return &FilesController{app: app}
}

func (fc *FilesController) Upload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(5 << 20) // 5 MB
	if err != nil {
		WriteJSONErr(w, &models.Error{Code: "invalid form data", Err: fmt.Errorf("failed to parse multipart form: %w", err)})
		return
	}

	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)

	var files []*models.File

	for _, handler := range r.MultipartForm.File["files"] {
		data, err := handler.Open()
		if err != nil {
			WriteHTMLErr(w, fmt.Errorf("failed opening file: %w", err))
			return
		}

		defer data.Close()

		file, err := fc.app.UploadFile(ctx, &commands.UploadFileRequest{
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
	pages.FileGrid(files, models.FileURL).Render(ctx, w)

}

func (fc *FilesController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")

	if err := fc.app.DeleteFile(ctx, config.OrganizationID, filename); err != nil {
		WriteHTMLErr(w, err)
		return
	}

	w.Header().Set("HX-Trigger", "deleted, render")
	w.WriteHeader(http.StatusOK)
}

func (fc *FilesController) Download(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)

	var err error
	var req queries.GetFileDataRequest

	req.Name = r.PathValue("filename")
	req.OrganizationID = config.OrganizationID

	if variant := r.URL.Query().Get("variant"); variant != "" {
		if req.Variant, err = strconv.Atoi(variant); err != nil {
			WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid quality: %w", err)))
			return
		}
	}

	res, err := fc.app.GetFileData(r.Context(), &req)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed getting file: %w", err))
		return
	}

	w.Header().Set("HX-Trigger", "render")
	WriteFile(w, r, res)

}

func (fc *FilesController) Rename(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(middlewares.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")

	if err := r.ParseForm(); err != nil {
		WriteHTMLErr(w, models.DecodeError(err))
	}

	newName := r.PostFormValue("name")

	if err := fc.app.RenameFile(ctx, config.OrganizationID, filename, newName); err != nil {
		WriteHTMLErr(w, fmt.Errorf("failed renaming file: %w", err))
		return
	}

	w.Header().Set("HX-Trigger", "renamed, render")
	w.WriteHeader(http.StatusOK)
}
