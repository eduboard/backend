package main

import (
	"fmt"
	"github.com/eduboard/backend/http"
	"github.com/eduboard/backend/mongodb"
	"github.com/eduboard/backend/service"
	"os"
)

func main() {
	var mongoURL = "mongodb://localhost:27017"
	v, ok := os.LookupEnv("MONGO_URL")
	if ok {
		mongoURL = v
	}

	repository := mongodb.Initialize(mongoURL)

	uS := service.NewUserService(repository.UserRepository)
	cS := service.NewCourseService(repository.CourseRepository)

	server := http.AppServer{
		UserService:   uS,
		CourseService: cS,
	}

	fmt.Println("listening now")
	server.Run()
}
