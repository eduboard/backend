package http

import (
	"github.com/eduboard/backend"
	"log"
	"net/http"
	"time"
)

type AppServer struct {
	Host          string
	Static        string
	UserService   eduboard.UserService
	CourseService eduboard.CourseService
	httpServer    *http.Server
}

func (a *AppServer) initialize() {
	protected := a.authenticatedRoutes()
	public := a.publicRoutes()

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", NewAuthMiddleware(a.UserService, protected))
	mux.Handle("/api/", public)
	mux.Handle("/", http.FileServer(http.Dir(a.Static)))

	a.httpServer = &http.Server{
		Addr:           a.Host,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}
}

func (a *AppServer) Run() {
	a.initialize()
	log.Fatal(a.httpServer.ListenAndServe())
}
