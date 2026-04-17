package middlewares

import (
	"context"
	"net/http"
)

const PathContextKey = key("path")

func WithPath(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		nr := r.WithContext(context.WithValue(r.Context(), PathContextKey, path))
		next.ServeHTTP(w, nr)
	}
}
