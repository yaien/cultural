package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"github.com/yaien/cultural/internal/modules/configs/internal/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewWithUser(app *application.Application, store sessions.Store) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s, _ := store.Get(r, controllers.SessionName)
			id, ok := s.Values[controllers.UserIDKey].(string)
			if !ok || id == "" {
				redirect(s, w, r)
				return
			}

			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, "Invalid user ID in session", http.StatusInternalServerError)
				return
			}

			user, err := app.GetUserByID(r.Context(), oid)
			if err != nil {
				redirect(s, w, r)
				return
			}

			nr := r.WithContext(context.WithValue(r.Context(), models.UserContextKey, user))
			next.ServeHTTP(w, nr)
		}
	}
}

func redirect(s *sessions.Session, w http.ResponseWriter, r *http.Request) {
	s.AddFlash(r.URL.Path, controllers.RedirectKey)

	if err := s.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/auth/google/login", http.StatusPermanentRedirect)
}
