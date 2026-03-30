package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func WithCache(next http.Handler) http.HandlerFunc {
	now := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		etag := fmt.Sprintf(`W/"%d"`, now.Unix())
		modified := now.Format(http.TimeFormat)

		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
		w.Header().Set("Etag", etag)
		w.Header().Set("Last-Modified", modified)

		if match := r.Header.Get("If-None-Match"); match != "" {
			if match == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		if since := r.Header.Get("If-Modified-Since"); since != "" {
			if t, err := time.Parse(http.TimeFormat, since); err == nil && now.Before(t.Add(1*time.Second)) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
