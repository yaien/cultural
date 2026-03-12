package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func WriteJSONSuccess(w http.ResponseWriter) {
	WriteJSON(w, map[string]any{"success": true})
}

func WriteJSONErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var e *models.Error

	switch {
	case errors.As(err, &e):
		status := http.StatusBadRequest
		if e.HTTPStatus > 0 {
			status = e.HTTPStatus
		}

		w.WriteHeader(status)

		err = json.NewEncoder(w).Encode(map[string]string{"error": e.Code, "message": e.Error()})
	default:
		slog.Error("Internal server error", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	}

	if err != nil {
		slog.Warn("Failed to write JSON error response", "err", err.Error())
	}
}

func WriteHTMLErr(w http.ResponseWriter, err error) {
	var e *models.Error

	switch {
	case errors.As(err, &e):
		status := http.StatusBadRequest
		if e.HTTPStatus > 0 {
			status = e.HTTPStatus
		}

		w.WriteHeader(status)
		_, err = fmt.Fprintf(w, "<h1>Error: %s - %s</h1>", e.Code, e.Error())
	default:
		slog.Error("Internal server error", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
	}

	if err != nil {
		slog.Warn("Failed to write HTML error response", "err", err.Error())
	}
}

func WriteFile(w http.ResponseWriter, r *http.Request, res *queries.GetFileDataResponse) {
	defer func() {
		if err := res.Data.Close(); err != nil {
			slog.Warn("error closing file data", "err", err)
		}
	}()

	w.Header().Set("Content-Type", res.ContentType)
	w.Header().Set("Content-Length", fmt.Sprint(res.Size))
	w.Header().Set("Cache-Control", "public, max-age=0")
	w.Header().Set("ETag", res.ID.Hex())
	w.Header().Set("Last-Modified", res.UpdatedAt.Format(http.TimeFormat))

	if match := r.Header.Get("If-None-Match"); match != "" && match == res.ID.Hex() {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	if modified := r.Header.Get("If-Modified-Since"); modified != "" {
		if t, err := http.ParseTime(modified); err == nil && res.UpdatedAt.Equal(t) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	if _, err := io.Copy(w, res.Data); err != nil {
		slog.Warn("error downloading the file", "err", err)
		return
	}
}
