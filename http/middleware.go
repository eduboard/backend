package http

import (
	"fmt"
	"github.com/eduboard/backend"
	"net/http"
)

// Chain takes a final http.Handler and a list of Middlewares and builds a call chain such that
// an incoming request passes all Middlwares in the order they were appended and finally reaches final.
func Chain(final http.Handler, m ...func(handler http.Handler) http.Handler) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		final = m[i](final)
	}
	return final
}

func NewAuthMiddleware(provider eduboard.UserAuthenticationProvider) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("sessionID")
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			sessionID := cookie.Value
			err, userID := provider.CheckAuthentication(sessionID)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			r.Header.Set("userID", userID)
			next.ServeHTTP(w, r)
		})
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}
