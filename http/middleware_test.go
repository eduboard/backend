package http

import (
	"errors"
	"github.com/eduboard/backend/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppServer_NewAuthMiddleware(t *testing.T) {
	var as = &mock.UserAuthenticationProvider{
		CheckAuthenticationFn: func(sessionID string) (err error, ok bool) {
			if sessionID == "" {
				return errors.New("empty sessionID"), false
			}
			if sessionID == "invalid" {
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

		var testHandler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
			assert.True(t, v.enter, "handler should not have been entered")
			handlerEntered = true
		}

		t.Run(v.name, func(t *testing.T) {
			handler := NewAuthMiddleware(as)(testHandler)
			if v.enter {
				assert.True(t, handlerEntered, "handler was not entered")
			}

			req := httptest.NewRequest("", "/", nil)
			if v.cookie {
				req.AddCookie(&http.Cookie{Name: "sessionID", Value: v.input})
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, v.status, rr.Code, "status code does not match")

			assert.True(t, as.CheckAuthenticationFnInvoked, "authentication was not actually checked")
		})
	}
}