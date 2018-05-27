package eduboard

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
	Courses  []string      `json:"courses" bson:"courses"`
}

type UserRepository interface {
	Store(user *User) error
	Find(id string) (error, *User)
}

type UserService interface {
	CreateUser(u *User) (error, *User)
	Login(username string, password string) (error, *User)
}
