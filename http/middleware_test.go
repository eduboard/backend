package http

import (
	"errors"
	"github.com/eduboard/backend/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"log"
)

func TestChain(t *testing.T) {
	var final, oneCalled, twoCalled, threeCalled = &mock.Check{}, &mock.Check{}, &mock.Check{}, &mock.Check{}
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, oneCalled.Passed, "first wrapped function not called")
		assert.True(t, twoCalled.Passed, "second wrapped function not called")
		assert.True(t, threeCalled.Passed, "third wrapped function not called")
		final.Passed = true
		n, err := w.Write([]byte("ok"))
		assert.Equal(t, 2, n)
		assert.Nil(t, err)
	})

	handlers := mock.GenerateCheckedMiddlewares(oneCalled, twoCalled, threeCalled)
	c := Chain(finalHandler, handlers...)
	req := httptest.NewRequest("", "/", nil)
	rr := httptest.NewRecorder()
	c.ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)

	assert.True(t, final.Passed, "inner wrapped function not called")
	assert.Equal(t, "ok", string(resBody), "response not correct")
}

func TestAppServer_NewAuthMiddleware(t *testing.T) {
	var as = &mock.UserAuthenticationProvider{
		CheckAuthenticationFn: func(sessionID string) (err error, id string) {
			if sessionID == "" {
				return errors.New("empty sessionID"), ""
			}
			if sessionID == "invalid" {
				return errors.New("not found"), ""
			}
			return nil, "1"
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
		handlerEntered := false

		var testHandler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
			assert.True(t, v.enter, "handler should not have been entered")
			handlerEntered = true
		}

		t.Run(v.name, func(t *testing.T) {
			handler := NewAuthMiddleware(as)(testHandler)

			req := httptest.NewRequest("", "/", nil)
			if v.cookie {
				req.AddCookie(&http.Cookie{Name: "sessionID", Value: v.input})
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, v.status, rr.Code, "status code does not match")
			assert.True(t, as.CheckAuthenticationFnInvoked, "authentication was not actually checked")
			assert.Equal(t, v.enter, handlerEntered, "handler was not called as expected")
		})
	}
}

func TestCORS(t *testing.T) {
	var testCases = []struct {
		name   string
		method string
		enter  bool
	}{
		{"GET", "GET", true},
		{"PUT", "PUT", true},
		{"POST", "POST", true},
		{"OPTIONS", "OPTIONS", false},
	}

	var expectedHeader = map[string]string{
		"Access-Control-Allow-Origin":  "http://localhost:8080",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding",
		"Access-Control-Allow-Credentials": "true",
	}

	for _, v := range testCases {
		handlerEntered := false

		var testHandler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
			assert.True(t, v.enter, "handler should not have been entered")
			handlerEntered = true
		}

		t.Run(v.name, func(t *testing.T) {
			handler := CORS(testHandler)
			req := httptest.NewRequest(v.method, "/", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			res := rr.Result()
			h := res.Header

			for k, v := range expectedHeader {
				assert.Equalf(t, v, h.Get(k), "header %s does not match expected value", k)
			}
			assert.Equal(t, http.StatusOK, res.StatusCode, "http status not 200")
			assert.Equal(t, v.enter, handlerEntered, "next handler was not called as expected")
		})
	}
}

func TestLogger(t *testing.T) {
	var testCases = []struct {
		name   string
		method string
		path   string
	}{
		{"login", "POST", "/api/login"},
		{"getCourses", "GET", "/api/v1/courses"},
		{"CORS", "OPTIONS", "/"},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			handlerEntered := false

			var testHandler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
				handlerEntered = true
			}

			b := &bytes.Buffer{}
			l := log.New(b, "", 0)
			handler := Logger(l)
			req := httptest.NewRequest(v.method, v.path, nil)
			rr := httptest.NewRecorder()

			handler(testHandler).ServeHTTP(rr, req)
			assert.True(t, handlerEntered, "inner handler was not entered")
			assert.Contains(t, b.String(), v.path, "log does not contain path")
			assert.Contains(t, b.String(), v.method, "log does not contain method")
		})
	}
}
