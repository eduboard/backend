package http

import (
	"github.com/eduboard/backend"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func NewAuthMiddleware(provider eduboard.UserAuthenticationProvider, nextHandler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

		nextHandler(w, r, p)
	}
}
