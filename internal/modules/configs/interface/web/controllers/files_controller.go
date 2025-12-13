package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/application/commands"
	"github.com/yaien/cultural/internal/modules/configs/models"
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
		http.Error(w, "error retrieving the file", http.StatusBadRequest)
		return
	}

	defer data.Close()

	file, err := fc.app.UploadFile(r.Context(), &commands.UploadFileRequest{
		Name:           handler.Filename,
		Size:           handler.Size,
		MimeType:       handler.Header.Get("Content-Type"),
		OrganizationID: config.OrganizationID,
		Data:           data,
	})

	if err != nil {
		http.Error(w, "error uploading the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(file)
}

func (fc *FileController) Delete(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	filename := r.PathValue("filename")

	err := fc.app.DeleteFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		http.Error(w, "error deleting the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fc *FileController) List(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	files, err := fc.app.GetFiles(r.Context(), config.OrganizationID)
	if err != nil {
		http.Error(w, "error listing files", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(files)
}

func (fc *FileController) Get(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")

	file, _, err := fc.app.GetFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		http.Error(w, "error retrieving the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(file)
}

func (fc *FileController) Download(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")
	file, data, err := fc.app.GetFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		http.Error(w, "error retrieving the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Length", fmt.Sprint(file.Size))
	w.WriteHeader(http.StatusOK)

	_, err = bufio.NewWriter(w).ReadFrom(data)
	if err != nil {
		http.Error(w, "error downloading the file", http.StatusInternalServerError)
		return
	}
}

func (fc *FileController) Rename(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")

	var input struct {
		NewName string `json:"newName"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = fc.app.RenameFile(r.Context(), config.OrganizationID, filename, input.NewName)
	if err != nil {
		http.Error(w, "error renaming the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
