package http

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (a *AppServer) registerUserHandler() httprouter.Handle {
	type request struct{}
	type response struct{}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Do something with a.UserService
	}
}

func (a *AppServer) getUserHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Do something with a.UserService
	}
}

func (a *AppServer) getAllCoursesHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Do something with a.CourseService
	}
}

func (a *AppServer) getCoursesHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// Do something with a.CourseService
	}
}
