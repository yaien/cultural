package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/yaien/cultural/internal/application/storage"
	"github.com/yaien/cultural/internal/lib/coderror"
	"github.com/yaien/cultural/internal/lib/primitive"
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

	var e *coderror.Error

	switch {
	case errors.As(primitive.Error(err), &e):
		w.WriteHeader(e.HTTPStatus())

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
	var e *coderror.Error

	switch {
	case errors.As(primitive.Error(err), &e):

		w.WriteHeader(e.HTTPStatus())
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

func WriteFile(w http.ResponseWriter, r *http.Request, d *storage.Download) {
	defer func() {
		if err := d.Data.Close(); err != nil {
			slog.Warn("error closing file data", "err", err)
		}
	}()

	w.Header().Set("Content-Type", d.ContentType)
	w.Header().Set("Content-Length", fmt.Sprint(d.Size))
	w.Header().Set("Cache-Control", "public, max-age=0")
	w.Header().Set("ETag", d.ID)
	w.Header().Set("Last-Modified", d.UpdatedAt.Format(http.TimeFormat))

	if match := r.Header.Get("If-None-Match"); match != "" && match == d.ID {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	if modified := r.Header.Get("If-Modified-Since"); modified != "" {
		if t, err := http.ParseTime(modified); err == nil && d.UpdatedAt.Equal(t) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	if _, err := io.Copy(w, d.Data); err != nil {
		slog.Warn("error downloading the file", "err", err)
		return
	}
}
