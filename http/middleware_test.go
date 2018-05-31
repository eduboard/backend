package http

import (
	"errors"
	"github.com/eduboard/backend/mock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppServer_NewAuthMiddleware(t *testing.T) {
	var as = &mock.UserAuthenticationProvider{
		CheckAuthenticationFn: func(sessionId string) (err error, ok bool) {
			if sessionId == "" {
				return errors.New("empty sessionId"), false
			}
			if sessionId == "invalid" {
				return errors.New("not found"), false
			}
			return nil, true
		},
	}

	var testCases = []struct {
		name   string
		enter  bool
		status int
		input  string
		cookie bool
	}{
		{"empty sessionid", false, 403, "", true},
		{"invalid sessionid", false, 403, "invalid", true},
		{"no cookie", false, 403, "whatever", false},
		{"success", true, 200, "valid", true},
	}

	for _, v := range testCases {
		handlerEntered := true
		t.Run(v.name, func(t *testing.T) {
			handler := NewAuthMiddleware(as, func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
				assert.True(t, v.enter, "handler should not have been entered")
				handlerEntered = true
			})
			if v.enter {
				assert.True(t, handlerEntered, "handler was not entered")
			}

			router := httprouter.New()
			router.GET("/", handler)
			req, _ := http.NewRequest("GET", "/", nil)

			if v.cookie {
				req.AddCookie(&http.Cookie{Name: "sessionID", Value: v.input})
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, v.status, rr.Code, "status code does not match")

			assert.True(t, as.CheckAuthenticationFnInvoked, "authentication was not actually checked")
		})
	}
}
