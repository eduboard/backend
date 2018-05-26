package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func main() {

	router := httprouter.New()
	initializeRoutes(router)

	server := &http.Server{
		Addr: "localhost:8080",
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler : router,
	}

	server.ListenAndServe()
}
