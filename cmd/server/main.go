package main

import (
	"fmt"
	"github.com/eduboard/backend/http"
	"github.com/eduboard/backend/mongodb"
	"github.com/eduboard/backend/service"
)

func main() {
	repository := mongodb.Initialize()

	uS := service.NewUserService(repository.UserRepository)
	cS := service.NewCourseService(repository.CourseRepository)

	server := http.AppServer{
		UserService:   uS,
		CourseService: cS,
	}

	fmt.Println("listening now")
	server.Run()
}
