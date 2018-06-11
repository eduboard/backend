package http

import (
	"github.com/julienschmidt/httprouter"
)

func (a *AppServer) authenticatedRoutes() *httprouter.Router {
	router := httprouter.New()

	// User
	router.GET("/api/v1/user/:id", a.getUserHandler())

	// Courses
	router.GET("/api/v1/courses/:id", a.getCourseHandler())
	router.GET("/api/v1/courses", a.getAllCoursesHandler())
	return router
}

func (a *AppServer) publicRoutes() *httprouter.Router {
	router := httprouter.New()

	// Registration
	router.POST("/api/register", a.registerUserHandler())
	router.POST("/api/login", a.loginUserHandler())
	router.POST("/api/logout", a.logoutUserHandler())
	return router
}
