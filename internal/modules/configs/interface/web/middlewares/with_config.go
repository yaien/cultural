package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

func NewWithConfig(app *application.Application) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
				host = forwardedHost
			}

			slog.Debug(
				"Request Received",
				"host", r.Host,
				"forwarded-host", r.Header.Get("X-Forwarded-Host"),
				"path", r.URL.Path,
				"method", r.Method,
				"url", r.URL.String(),
			)

			if host == "" {
				http.Error(w, "Missing host header", http.StatusBadRequest)
				return
			}

			config, err := app.GetConfigByHost(r.Context(), host)
			if err != nil {
				http.Error(w, "Failed to get config", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), models.ConfigContextKey, config)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
