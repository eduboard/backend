package http

import (
	"github.com/eduboard/backend"
	"log"
	"net/http"
	"time"
)

type AppServer struct {
	UserService   eduboard.UserService
	CourseService eduboard.CourseService
	httpServer    *http.Server
}

func (a *AppServer) initialize() {
	router := a.authenticatedRoutes()

	a.httpServer = &http.Server{
		Addr:           ":8080",
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}
}

func (a *AppServer) Run() {
	a.initialize()
	log.Fatal(a.httpServer.ListenAndServe())
}
