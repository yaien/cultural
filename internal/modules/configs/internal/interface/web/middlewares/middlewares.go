package middlewares

import "net/http"

type Middlewares struct {
	WithConfig func(next http.Handler) http.HandlerFunc
	WithUser   func(next http.Handler) http.HandlerFunc
	WithRole   func(next http.Handler) http.HandlerFunc
}
