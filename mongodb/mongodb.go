package mongodb

import (
	"gopkg.in/mgo.v2"
	"log"
)

type DBConfig struct {
	URL      string
	Database string
	Username string
	Password string
}

type Repository struct {
	session               *mgo.Session
	UserRepository        *UserRepository
	CourseRepository      *CourseRepository
	CourseEntryRepository *CourseEntryRepository
}

func Initialize(c DBConfig) *Repository {
	config := &mgo.DialInfo{
		Addrs:    []string{c.URL},
		Database: c.Database,
		Username: c.Username,
		Password: c.Password,
	}

	session, err := mgo.DialWithInfo(config)
	if err != nil {
		log.Fatalf("error connecting to mongoDB: %v", err)
	}

	db := session.DB(config.Database)
	return &Repository{
		session:               session,
		UserRepository:        newUserRepository(db),
		CourseRepository:      newCourseRepository(db),
		CourseEntryRepository: newCourseEntryRepository(db),
	}
}
