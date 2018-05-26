package user

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

type User struct {
	Username string `json:"username" bson:"username"`
	Id       string `json:"id" bson:"_id"`
}

type Repository struct {
	Session *mgo.Session
	C       *mgo.Collection
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) InitializeDB() {
	type DB struct {
		url        string
		database   string
		collection string
		user       string
		password   string
	}

	dbURL, ok := os.LookupEnv("MONGO_URL")
	if !ok {
		dbURL = "localhost:27017"
	}

	db := DB{
		url:        dbURL,
		database:   "eduboard",
		collection: "users",
		user:       "",
		password:   "",
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{db.url},
		Database: db.database,
		Username: db.user,
		Password: db.password,
	})
	if err != nil {
		log.Fatalf("error connecting to mongoDB %v", err)
	}

	r.Session = session
	r.C = session.DB(db.database).C(db.collection)
}

func (r *Repository) FindUser(id string) (err error, u User) {

	result := User{}

	err = r.C.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return err, User{}
	}

	return nil, result

}
