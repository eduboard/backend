package main

import (
	"fmt"
	"github.com/eduboard/backend/config"
	"github.com/eduboard/backend/http"
	"github.com/eduboard/backend/mongodb"
	"github.com/eduboard/backend/service"
)

func main() {

	c := config.GetConfig()
	mongoConfig := mongodb.DBConfig{
		URL:      fmt.Sprintf("%s:%s", c.MongoHost, c.MongoPort),
		Database: c.MongoDB,
		Username: c.MongoUser,
		Password: c.MongoPass,
	}

	repository := mongodb.Initialize(mongoConfig)

	uS := service.NewUserService(repository.UserRepository)
	cS := service.NewCourseService(repository.CourseRepository)

	server := http.AppServer{
		UserService:   uS,
		CourseService: cS,
	}

	fmt.Println("listening now")
	server.Run()
}
