package middlewares

import "net/http"

type Middleware func(next http.Handler) http.HandlerFunc

type Middlewares struct {
	WithConfig Middleware
	WithUser   Middleware
	WithRole   Middleware
	WithCache  Middleware
}
