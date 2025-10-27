package middlewares

import (
	"context"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/models"
)

type key string

const ConfigContextKey = key("config")

func NewWithConfig(app *application.Application) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			host := r.Host

			if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
				host = forwardedHost
			}

			if host == "" {
				http.Error(w, "Missing host header", http.StatusBadRequest)
				return
			}

			config, err := app.GetConfigByHost(r.Context(), host)
			if err != nil {
				http.Error(w, "Failed to get config", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), ConfigContextKey, config)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func GetConfig(ctx context.Context) *models.Config {
	return ctx.Value(ConfigContextKey).(*models.Config)
}
