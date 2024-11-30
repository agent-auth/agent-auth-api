package middlewares

import (
	"net/http"
)

// Middlewares interfaces
type Middlewares interface {
	Logger() func(h http.Handler) http.Handler
}
