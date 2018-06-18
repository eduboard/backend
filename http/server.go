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
	Logger        *log.Logger
	UserService   eduboard.UserService
	CourseService eduboard.CourseService
	httpServer    *http.Server
}

func (a *AppServer) initialize() {
	protected := a.authenticatedRoutes()
	public := a.publicRoutes()

	privateChain := Chain(protected, Logger(a.Logger), CORS, NewAuthMiddleware(a.UserService))
	publicChain := Chain(public, Logger(a.Logger), CORS)
	staticChain := Chain(http.FileServer(http.Dir(a.Static)), Logger(a.Logger), CORS)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", privateChain)
	mux.Handle("/api/", publicChain)
	mux.Handle("/", staticChain)

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
