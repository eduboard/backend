package http

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (a *AppServer) initializeRoutes(router *httprouter.Router) {
	router.GET("/", a.serveFilesHandler())

	// Registration
	router.POST("/api/v1/register", a.registerUserHandler())
	router.POST("/api/v1/login", a.loginUserHandler())
	router.POST("/api/v1/logout", a.logoutUserHandler())

	// User
	router.GET("/api/v1/user/:id", a.getUserHandler())

	// Courses
	router.GET("/api/v1/courses/:id", a.getCourseHandler())
	router.GET("/api/v1/courses/", NewAuthMiddleware(a.UserService, a.getAllCoursesHandler()))
}

func (a *AppServer) serveFilesHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "hello world")
	}
}
