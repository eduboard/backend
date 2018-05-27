package http

import (
	"github.com/eduboard/backend"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type AppServer struct {
	UserService   eduboard.UserService
	CourseService eduboard.CourseService
	httpServer    *http.Server
}

func (a *AppServer) initialize() {
	router := httprouter.New()
	a.httpServer = &http.Server{
		Addr:           "localhost:8080",
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}

	a.initializeRoutes(router)
}

func (a *AppServer) Run() {
	a.initialize()
	a.httpServer.ListenAndServe()
}
