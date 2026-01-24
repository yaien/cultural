package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/modules/configs/application"
	"github.com/yaien/cultural/internal/modules/configs/interface/web/controllers"
	"github.com/yaien/cultural/internal/modules/configs/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewWithUser(app *application.Application, store sessions.Store) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, controllers.SessionName)
			id, ok := session.Values[controllers.UserIDKey].(string)
			if !ok || id == "" {
				session.AddFlash(r.URL.Path, controllers.RedirectKey)
				_ = session.Save(r, w)
				http.Redirect(w, r, "/auth/google/login", http.StatusTemporaryRedirect)
				return
			}

			oid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				http.Error(w, "Invalid user ID in session", http.StatusInternalServerError)
				return
			}

			user, err := app.GetUserByID(r.Context(), oid)
			if err != nil {
				session.AddFlash(r.URL.Path, controllers.RedirectKey)
				_ = session.Save(r, w)
				http.Redirect(w, r, "/auth/google/login", http.StatusTemporaryRedirect)
				return
			}

			nr := r.WithContext(context.WithValue(r.Context(), models.UserContextKey, user))
			next.ServeHTTP(w, nr)
		}
	}
}
