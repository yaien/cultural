package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/modules/configs/internal/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type key string

const (
	SessionKey        = "session"
	UserIDKey         = "user_id"
	RedirectKey       = "redirect"
	UserContextKey    = key("user")
	SessionContextKey = key("session")
)

func NewWithUser(app *application.Application, store sessions.Store) func(next http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s, _ := store.Get(r, SessionKey)
			id, ok := s.Values[UserIDKey].(string)
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
				http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserContextKey, user)
			ctx = context.WithValue(ctx, SessionContextKey, s)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func redirect(_ *sessions.Session, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, "/auth/google/login?redirect="+r.URL.Path, http.StatusPermanentRedirect)
}
