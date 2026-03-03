package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const RoleContextKey = key("role")

func NewWithRole(app *application.Application, store sessions.Store) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, SessionKey)
			id, ok := session.Values[UserIDKey].(string)
			if !ok || id == "" {
				redirect(session, w, r)
				return
			}

			userID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, "Invalid user ID in session", http.StatusInternalServerError)
				return
			}

			config, ok := r.Context().Value(ConfigContextKey).(*models.Config)
			if !ok || config == nil {
				http.Error(w, "Config not found in context", http.StatusInternalServerError)
				return
			}

			role, err := app.GetRole(r.Context(), userID, config.OrganizationID)
			if err != nil {
				http.Error(w, "Failed to get user role", http.StatusInternalServerError)
				return
			}

			nr := r.WithContext(context.WithValue(r.Context(), RoleContextKey, role))
			nr = nr.WithContext(context.WithValue(nr.Context(), SessionContextKey, session))
			next.ServeHTTP(w, nr)
		}
	}
}
