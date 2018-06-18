package main

import (
	"fmt"
	"github.com/eduboard/backend/config"
	"github.com/eduboard/backend/http"
	"github.com/eduboard/backend/mongodb"
	"github.com/eduboard/backend/service/courseEntryService"
	"github.com/eduboard/backend/service/courseService"
	"github.com/eduboard/backend/service/userService"
	"io"
	"log"
	"os"
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

	var logDst io.Writer
	if c.LogFile == "" {
		logDst = os.Stdout
	} else {
		file, err := os.OpenFile(c.LogFile, os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("error opening logfile %s: %v", c.LogFile, err)
		}
		defer file.Close()
		logDst = file
	}

	server := http.AppServer{
		Host:                  c.Host,
		Static:                c.StaticDir,
		Logger:                log.New(logDst, "", log.LstdFlags),
		UserService:           userService.New(repository.UserRepository),
		CourseService:         courseService.New(repository.CourseRepository),
		CourseEntryService:    courseEntryService.New(repository.CourseEntryRepository),
		CourseRepository:      repository.CourseRepository,
		CourseEntryRepository: repository.CourseEntryRepository,
	}

	server.Logger.Printf("Server listening on %s", c.Host)
	server.Run()
}
