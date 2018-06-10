package http

import (
	"github.com/julienschmidt/httprouter"
)

func (a *AppServer) authenticatedRoutes() *httprouter.Router {
	router := httprouter.New()

	// User
	router.GET("/api/v1/users/:id", a.GetUserHandler())
	router.GET("/api/v1/me", a.GetMeHandler())

	// Courses
	router.GET("/api/v1/courses/:id", a.getCourseHandler())
	router.GET("/api/v1/courses", a.getAllCoursesHandler())

	return router
}

func (a *AppServer) publicRoutes() *httprouter.Router {
	router := httprouter.New()

	// Registration
	router.POST("/api/register", a.RegisterUserHandler())
	router.POST("/api/login", a.LoginUserHandler())
	router.POST("/api/logout", a.LogoutUserHandler())
	return router
}
