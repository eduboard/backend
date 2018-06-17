package http

import (
	"errors"
	"github.com/eduboard/backend"
	"github.com/eduboard/backend/mock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"log"
	"os"
)

func TestAppServer_RegisterUserHandler(t *testing.T) {
	mockService := mock.UserService{}
	mockService.CreateUserFn = func(u *eduboard.User, password string) (error, *eduboard.User) {
		if len(password) < 8 {
			return errors.New("too short"), u
		}

		u.ID = bson.ObjectIdHex("5b1d24e72c5b292fe0d6ee55")
		return nil, u
	}
	appServer := AppServer{UserService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name   string
		body   string
		status int
	}{
		{"no email", `{"password":"password"}`, 400},
		{"no password", `{"email":"e@mail.com"}`, 400},
		{"password too short", `{"email":"e@mail.com","password":"pass"}`, 500},
		{"malformed json", `{"email":"e@mail.com","password":"pass"`, 400},
		{"success", `{"email":"e@mail.com","password":"password"}`, 200},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.CreateUserFnInvoked = false
			r := httptest.NewRequest("POST", "/", strings.NewReader(v.body))
			rr := httptest.NewRecorder()

			appServer.RegisterUserHandler()(rr, r, httprouter.Params{})

			if v.status == 400 {
				assert.Equal(t, v.status, rr.Code, "bad response code")
				assert.False(t, mockService.CreateUserFnInvoked, "CreateUser was invoked when it should not")
				assert.Empty(t, rr.Body, "body should be empty")
				return
			}

			assert.True(t, mockService.CreateUserFnInvoked, "CreateUser was not invoked when it should")
			assert.Equal(t, v.status, rr.Code, "bad response code")

			if v.status == 200 {
				assert.NotEmptyf(t, rr.HeaderMap["Set-Cookie"], "cookie was not set on successful registration")
			}
		})
	}
}

func TestAppServer_LoginUserHandler(t *testing.T) {
	mockService := mock.UserService{}
	mockService.LoginFn = func(email string, password string) (error, *eduboard.User) {
		if password != "password" {
			return errors.New("bad login"), &eduboard.User{}
		}

		user := &eduboard.User{ID: bson.ObjectIdHex("5b1d24e72c5b292fe0d6ee55")}
		return nil, user
	}
	appServer := AppServer{UserService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name   string
		body   string
		status int
	}{
		{"no email", `{"password":"password"}`, 400},
		{"no password", `{"email":"e@mail.com"}`, 400},
		{"bad password", `{"email":"e@mail.com","password":"pass"}`, 401},
		{"malformed json", `{"email":"e@mail.com","password":"pass"`, 400},
		{"success", `{"email":"e@mail.com","password":"password"}`, 200},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.LoginFnInvoked = false
			r := httptest.NewRequest("POST", "/", strings.NewReader(v.body))
			rr := httptest.NewRecorder()

			appServer.LoginUserHandler()(rr, r, httprouter.Params{})

			if v.status == 400 {
				assert.Equal(t, v.status, rr.Code, "bad response code")
				assert.False(t, mockService.LoginFnInvoked, "Login was invoked when it should not")
				assert.Empty(t, rr.Body, "body should be empty")
				return
			}

			assert.True(t, mockService.LoginFnInvoked, "Login was not invoked when it should")
			assert.Equal(t, v.status, rr.Code, "bad response code")

			if v.status == 200 {
				assert.NotEmptyf(t, rr.HeaderMap["Set-Cookie"], "cookie was not set on successful registration")
			}
		})
	}
}

func TestAppServer_LogoutUserHandler(t *testing.T) {
	mockService := mock.UserService{}
	mockService.LogoutFn = func(sessionID string) error {
		if sessionID != "session" {
			return errors.New("not found")
		}
		return nil
	}
	appServer := AppServer{UserService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name      string
		sessionID string
		status    int
	}{
		{"no sessionID", "", 200},
		{"bad sessionID", "535510N", 500},
		{"success", "session", 200},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.LogoutFnInvoked = false
			r := httptest.NewRequest("POST", "/", nil)
			c := http.Cookie{Name: "sessionID", Value: v.sessionID}
			r.AddCookie(&c)
			rr := httptest.NewRecorder()

			appServer.LogoutUserHandler()(rr, r, httprouter.Params{})

			assert.Equal(t, v.status, rr.Code, "bad response code")
			assert.Empty(t, rr.Body, "body should be empty")
			if v.sessionID == "" {
				assert.False(t, mockService.LogoutFnInvoked, "Logout was invoked when it should not")
				assert.Equal(t, v.status, rr.Code, "bad response code")
				return
			}
			assert.True(t, mockService.LogoutFnInvoked, "Logout was not invoked when it should")

			if v.status == 200 {
				assert.NotEmptyf(t, rr.HeaderMap["Set-Cookie"], "cookie was not set on successful registration")
			}
		})
	}
}

func TestAppServer_GetUserHandler(t *testing.T) {
	mockService := mock.UserService{}
	mockService.GetUserFn = func(id string) (error, *eduboard.User) {
		if id != "userId" {
			return errors.New("not found"), &eduboard.User{}
		}
		return nil, &eduboard.User{Name: "name"}
	}
	appServer := AppServer{UserService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name   string
		userID string
		status int
	}{
		{"no userID", "", 400},
		{"bad userID", "someId", 404},
		{"success", "userId", 200},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.GetUserFnInvoked = false
			r := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			appServer.GetUserHandler()(rr, r, httprouter.Params{httprouter.Param{"id", v.userID}})

			assert.Equal(t, v.status, rr.Code, "bad response code")
			if v.userID == "" {
				assert.False(t, mockService.GetUserFnInvoked, "GetUser was invoked when it should not")
				assert.Equal(t, v.status, rr.Code, "bad response code")
				assert.Empty(t, rr.Body, "body should not be empty")
				return
			}
			assert.True(t, mockService.GetUserFnInvoked, "GetUser was not invoked when it should")
			if v.status == 200 {
				assert.NotEmptyf(t, rr.Body, "body should not be empty")
			}
		})
	}
}

func TestAppServer_GetMeHandler(t *testing.T) {
	mockService := mock.UserService{}
	mockService.GetUserFn = func(id string) (error, *eduboard.User) {
		if id != "userId" {
			return errors.New("not found"), &eduboard.User{}
		}
		return nil, &eduboard.User{Name: "name"}
	}
	appServer := AppServer{UserService: &mockService, Logger: log.New(os.Stdout, "", 0)}

	var testCases = []struct {
		name   string
		userID string
		status int
	}{
		{"no userID", "", 400},
		{"bad userID", "someId", 404},
		{"success", "userId", 200},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			mockService.GetUserFnInvoked = false
			r := httptest.NewRequest("GET", "/", nil)
			r.Header.Set("userID", v.userID)
			rr := httptest.NewRecorder()

			appServer.GetMeHandler()(rr, r, httprouter.Params{})

			assert.Equal(t, v.status, rr.Code, "bad response code")
			if v.userID == "" {
				assert.False(t, mockService.GetUserFnInvoked, "GetUser was invoked when it should not")
				assert.Equal(t, v.status, rr.Code, "bad response code")
				assert.Empty(t, rr.Body, "body should not be empty")
				return
			}
			assert.True(t, mockService.GetUserFnInvoked, "GetUser was not invoked when it should")
			if v.status == 200 {
				assert.NotEmptyf(t, rr.Body, "body should not be empty")
			}
		})
	}
}
