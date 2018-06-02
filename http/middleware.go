package http

import (
	"github.com/eduboard/backend"
	"net/http"
)

func NewAuthMiddleware(provider eduboard.UserAuthenticationProvider, nextHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		accessToken := cookie.Value
		err, ok := provider.CheckAuthentication(accessToken)
		if err != nil || !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		nextHandler.ServeHTTP(w, r)
	}
}
