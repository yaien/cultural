package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/commands"
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
		MimeType:       handler.Header.Get("Content-Type"),
		OrganizationID: config.OrganizationID,
		Data:           data,
	})

	if err != nil {
		slog.Error("error uploading file", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(file)
}

func (fc *FileController) Delete(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)

	filename := r.PathValue("filename")

	err := fc.app.DeleteFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		slog.Error("error deleting file", "err", err)
		http.Error(w, "error deleting the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fc *FileController) List(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	files, err := fc.app.GetFiles(r.Context(), config.OrganizationID)
	if err != nil {
		slog.Error("error listing files", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "error listing files"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(files)
}

func (fc *FileController) Get(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")

	file, _, err := fc.app.GetFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		slog.Error("error retrieving file", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "error retrieving the file"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(file)
}

func (fc *FileController) Download(w http.ResponseWriter, r *http.Request) {
	config := r.Context().Value(models.ConfigContextKey).(*models.Config)
	filename := r.PathValue("filename")
	file, data, err := fc.app.GetFile(r.Context(), config.OrganizationID, filename)
	if err != nil {
		slog.Error("error retrieving file", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "error retrieving the file"})
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.Name+"\"")
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Length", fmt.Sprint(file.Size))

	_, err = bufio.NewWriter(w).ReadFrom(data)
	if err != nil {
		slog.Error("error downloading the file", "err", err)
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
		slog.Error("error parsing request body", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "invalid request body"})
		return
	}

	err = fc.app.RenameFile(r.Context(), config.OrganizationID, filename, input.NewName)
	if err != nil {
		slog.Error("error renaming the file", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "error renaming the file"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
