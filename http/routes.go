package http

import (
	"github.com/julienschmidt/httprouter"
)

func (a *AppServer) authenticatedRoutes() *httprouter.Router {
	router := httprouter.New()

	// User
	router.GET("/api/v1/users", a.GetAllUsersHandler())
	router.GET("/api/v1/users/:id", a.GetUserHandler())
	router.GET("/api/v1/users/:id/courses", a.GetMyCoursesHandler())
	router.GET("/api/v1/me", a.GetMeHandler())

	// Courses
	router.GET("/api/v1/courses/:courseID", a.GetCourseHandler())
	router.GET("/api/v1/courses/:courseID/users", a.GetMembersHandler())
	router.POST("/api/v1/courses/:courseID/users/subscribe", a.AddMembersHandler())
	router.POST("/api/v1/courses/:courseID/users/unsubscribe", a.RemoveMembersHandler())
	router.GET("/api/v1/courses", a.GetAllCoursesHandler())

	// CourseEntries
	router.POST("/api/v1/courses", a.CreateCourseHandler())
	router.POST("/api/v1/courses/:courseID/entries", a.PostCourseEntryHandler())
	router.PUT("/api/v1/courses/:courseID/entries/:entryID", a.PutCourseEntryHandler())
	router.DELETE("/api/v1/courses/:courseID/entries/:entryID", a.DeleteCourseEntryHandler())

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
