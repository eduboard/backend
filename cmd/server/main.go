package main

import (
	"fmt"
	"github.com/eduboard/backend/http"
	"github.com/eduboard/backend/mongodb"
	"github.com/eduboard/backend/service"
	"github.com/eduboard/backend/auth"
)

func main() {
	repository := mongodb.Initialize()
	authenticator := new(auth.UserAuthenticator)

	uS := service.NewUserService(repository.UserRepository, authenticator)
	cS := service.NewCourseService(repository.CourseRepository)

	server := http.AppServer{
		UserService:   uS,
		CourseService: cS,
	}

	fmt.Println("listening now")
	server.Run()
}
