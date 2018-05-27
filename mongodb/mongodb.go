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

func Initialize() *Repository {
	session, err := mgo.Dial("mongodb://localhost:27017")
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
