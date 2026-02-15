package middlewares

import (
	"context"
	"net/http"

	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
)

func NewWithRole(app *application.Application) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(models.UserContextKey).(*models.User)
			if !ok || user == nil {
				http.Error(w, "User not found in context", http.StatusUnauthorized)
				return
			}

			config, ok := r.Context().Value(models.ConfigContextKey).(*models.Config)
			if !ok || config == nil {
				http.Error(w, "Config not found in context", http.StatusInternalServerError)
				return
			}

			role, err := app.GetRole(r.Context(), user.ID, config.OrganizationID)
			if err != nil {
				http.Error(w, "Failed to get user role", http.StatusInternalServerError)
				return
			}

			nr := r.WithContext(context.WithValue(r.Context(), models.RoleContextKey, role))
			next.ServeHTTP(w, nr)
		}
	}
}
