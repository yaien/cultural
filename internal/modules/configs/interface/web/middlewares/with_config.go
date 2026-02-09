package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

func NewWithConfig(app *application.Application) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			if host, found := strings.CutPrefix(r.Host, "www."); found {
				http.Redirect(w, r, fmt.Sprintf("%s://%s%s", r.URL.Scheme, host, r.URL.Path), http.StatusMovedPermanently)
				return
			}

			slog.Debug(
				"Request Received",
				"host", host,
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
