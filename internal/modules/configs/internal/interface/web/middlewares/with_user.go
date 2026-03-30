package middlewares

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/yaien/cultural/internal/library/auth"
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

func NewWithUser(users *auth.Users, store sessions.Store) func(next http.Handler) http.HandlerFunc {
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

			ctx := r.Context()

			user, err := users.GetByID(ctx, oid)
			if err != nil {
				http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, UserContextKey, user)
			ctx = context.WithValue(ctx, SessionContextKey, s)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func redirect(_ *sessions.Session, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/google/login?redirect="+r.URL.Path, http.StatusTemporaryRedirect)
}
