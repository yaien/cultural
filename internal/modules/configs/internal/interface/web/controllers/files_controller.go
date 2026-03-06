package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/views/dashboard"
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

	templ.Join(
		pages.FileGrid(files, models.FileURL),
		dashboard.Toast("Archivos subido correctamente", dashboard.Primary),
	).Render(ctx, w)

}

func (fc *FilesController) Delete(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)

	filename := r.PathValue("filename")

	err := fc.app.DeleteFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		WriteJSONErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fc *FilesController) List(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)
	files, err := fc.app.GetFiles(r.Context(), config.OrganizationID)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed listing files: %w", err))
		return
	}

	WriteJSON(w, files)
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

	WriteFile(w, r, res)

}

func (fc *FilesController) Rename(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(middlewares.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")

	var input struct {
		NewName string `json:"newName"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		WriteJSONErr(w, models.DecodeError(err))
		return
	}

	err = fc.app.RenameFile(r.Context(), config.OrganizationID, filename, input.NewName)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed renaming file: %w", err))
		return
	}

	WriteJSONSuccess(w)
}
