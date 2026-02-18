package controllers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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
		return
	default:
		slog.Error("Internal server error", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
		fmt.Fprintf(w, "<h1>Error: %s</h1>", e.Code)
	default:
		slog.Error("Internal server error", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
	}
}

func WriteFile(w http.ResponseWriter, name, typ string, size int64, data io.ReadCloser) {
	defer data.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", name))
	w.Header().Set("Content-Type", typ)
	w.Header().Set("Content-Length", fmt.Sprint(size))

	_, err := bufio.NewWriter(w).ReadFrom(data)
	if err != nil {
		slog.Error("error downloading the file", "err", err)
		return
	}
}
