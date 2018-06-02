package mongodb

import (
	"gopkg.in/mgo.v2"
	"log"
)

type Repository struct {
	session          *mgo.Session
	UserRepository   *UserRepository
	CourseRepository *CourseRepository
}

func Initialize(mongoURL string ) *Repository {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Fatalf("error connecting to mongoDB: %v", err)
	}

	db := session.DB("eduboard")
	return &Repository{
		session:          session,
		UserRepository:   newUserRepository(db),
		CourseRepository: newCourseRepository(db),
	}
}
