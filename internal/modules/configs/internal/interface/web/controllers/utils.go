package controllers

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/modules/configs/internal/application/queries"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/middlewares"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

func init() {
	gob.Register(Toast{})
}

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
		_, err = fmt.Fprintf(w, "<h1>Error: %s</h1>", e.Code)
	default:
		slog.Error("Internal server error", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
	}

	if err != nil {
		slog.Warn("Failed to write HTML error response", "err", err.Error())
	}
}

func WriteFile(w http.ResponseWriter, r *http.Request, res *queries.GetFileResponse) {
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

type Toast struct {
	Message string
	Variant string
	Trigger string
}

const ToastKey = "toast"

// WriteSession adds the session to the request context and saves it in the response.
func WriteToast(w http.ResponseWriter, r *http.Request, toast Toast) {
	session, ok := r.Context().Value(middlewares.SessionContextKey).(*sessions.Session)
	if !ok {
		slog.Error("Failed to get session from context for toast message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.AddFlash(toast, ToastKey)
	if err := session.Save(r, w); err != nil {
		slog.Error("Failed to save session for toast message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	trigger := ToastKey
	if toast.Trigger != "" {
		trigger = toast.Trigger
	}

	w.Header().Set("HX-Trigger-After-Swap", trigger)
	w.WriteHeader(http.StatusOK)
}

// GetToast retrieves a toast message from the session in the request context.
func GetToast(w http.ResponseWriter, r *http.Request) (*Toast, bool, error) {
	session, ok := r.Context().Value(middlewares.SessionContextKey).(*sessions.Session)
	if !ok {
		return nil, false, fmt.Errorf("failed to get session from context for toast message")
	}

	toasts := session.Flashes(ToastKey)
	if len(toasts) == 0 {
		return nil, false, nil
	}

	store := session.Store()
	if store == nil {
		return nil, false, fmt.Errorf("session store is nil when retrieving toast message")
	}

	if err := session.Save(r, w); err != nil {
		return nil, false, fmt.Errorf("failed to save session after retrieving toast message: %w", err)
	}

	toast, ok := toasts[0].(Toast)
	if !ok {
		return nil, false, fmt.Errorf("invalid toast message type: %T", toasts[0])
	}

	return &toast, true, nil
}
