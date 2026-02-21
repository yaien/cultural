package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

type FileController struct {
	app *application.Application
}

func NewFileController(app *application.Application) *FileController {
	return &FileController{app: app}
}

func (fc *FileController) Upload(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	r.ParseMultipartForm(5 << 20) // 5 MB

	data, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("error retrieving the file", "err", err)
		http.Error(w, "error retrieving the file", http.StatusBadRequest)
		return
	}

	defer data.Close()

	file, err := fc.app.UploadFile(r.Context(), &commands.UploadFileRequest{
		Name:           handler.Filename,
		Size:           handler.Size,
		Type:           handler.Header.Get("Content-Type"),
		OrganizationID: config.OrganizationID,
		Data:           data,
	})

	if err != nil {
		WriteJSONErr(w, err)
		return
	}

	WriteJSON(w, file)
}

func (fc *FileController) Delete(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	filename := r.PathValue("filename")

	err := fc.app.DeleteFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		WriteJSONErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fc *FileController) List(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	files, err := fc.app.GetFiles(r.Context(), config.OrganizationID)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed listing files: %w", err))
		return
	}

	WriteJSON(w, files)
}

func (fc *FileController) Download(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	var err error
	var req queries.GetFileRequest

	req.Name = r.PathValue("filename")
	req.OrganizationID = config.OrganizationID

	if variant := r.URL.Query().Get("variant"); variant != "" {
		if req.Variant, err = strconv.Atoi(variant); err != nil {
			WriteJSONErr(w, models.DecodeError(fmt.Errorf("invalid quality: %w", err)))
			return
		}
	}

	res, err := fc.app.GetFile(r.Context(), &req)
	if err != nil {
		WriteJSONErr(w, fmt.Errorf("failed getting file: %w", err))
		return
	}

	WriteFile(w, res.Name, res.Type, res.Size, res.Data)

}

func (fc *FileController) Rename(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
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
