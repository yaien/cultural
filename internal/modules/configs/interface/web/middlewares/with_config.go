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
			if r.Header.Get("X-Forwarded-Host") != "" {
				host = r.Header.Get("X-Forwarded-Host")
			}

			scheme := r.Header.Get("X-Forwarded-Proto")
			if scheme == "" {
				if r.TLS != nil {
					scheme = "https"
				}
				scheme = "http"
			}

			if host, found := strings.CutPrefix(r.Host, "www."); found {
				http.Redirect(w, r, fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path), http.StatusMovedPermanently)
				return
			}

			slog.Debug(
				"Request",
				"host", host,
				"scheme", scheme,
				"path", r.URL.Path,
				"method", r.Method,
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
