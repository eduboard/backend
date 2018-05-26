package main

import (
	"github.com/eduboard/backend/controller/user"
	"github.com/julienschmidt/httprouter"
)

func initializeRoutes(router *httprouter.Router) {
	router.GET("/", user.RegisterUserHandler())
	// Registration
	//	router.POST("/api/v1/register")
	//	router.POST("/api/v1/login")
	//	router.POST("/api/v1/logout")

	// User
	router.GET("/api/v1/user/:id", user.GetUserHandler())

	// Courses
	//	router.GET("/api/v1/courses/")
	//		router.GET("/api/v1/courses/:id")*/
}
